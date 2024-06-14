import logging
import os
import re
from io import BytesIO
from pathlib import Path
from typing import List, Optional

import numpy as np
import torch
from fastapi import APIRouter
from pydantic import BaseModel
from scipy.io import wavfile

import ChatTTS
from app.common import global_info
from app.common.cached import info as cf
from app.utils import wav_utils
from app.utils.torch_utils import print_current_device_info

router = APIRouter()


class GenerateAudioRequest(BaseModel):
    text: str
    speed: int
    temperature: float
    top_P: float
    top_K: int
    refine_oral: int
    refine_laugh: int
    refine_break: int
    audio_seed_input: int
    text_seed_input: int
    refine_text_flag: bool
    customer_speaker_flag: bool
    customer_promot_flag: bool
    customer_promot: Optional[str] = ''
    customer_speaker_tensor: Optional[List] = None

    class Config:
        # 如果customer_speaker_tensor为None,则转换为空列表
        allow_population_by_field_name = True


class GenerateAudioResponse(BaseModel):
    output_text: str
    audio_name: str
    tensor: list

class GetAudioRequest(BaseModel):
    audio_path: str


class GetAudioResponse(BaseModel):
    # 音频文件
    audio_path: str


def create_chat(local_path=None):
    """创建ChatTTS实例"""
    if "chat" in cf:
        logging.info("cached chat")
        return cf['chat']
    local_model_path = os.getenv('MODEL_LOCAL_PATH', None)
    # 导入模型实例
    chat = ChatTTS.Chat()
    if local_model_path is None:
        chat.load_models()
    else:
        logging.info('local model path:', local_model_path)
        chat.load_models('local', local_path=local_model_path)
    cf['chat'] = chat
    return chat


def determine_seed(seed):
    """限定模型使用的种子值"""
    torch.manual_seed(seed)
    np.random.seed(seed)
    torch.cuda.manual_seed(seed)
    torch.backends.cudnn.deterministic = True
    torch.backends.cudnn.benchmark = False


def build_random_speaker(chat, request: GenerateAudioRequest):
    """
    生成随机音色
    :param chat:
    :param request:
    :return:
    """
    # 使用音色种子值创建音色
    audio_seed_input = request.audio_seed_input
    determine_seed(audio_seed_input)
    return chat.sample_random_speaker()


def build_code_text(chat, rand_spk, request: GenerateAudioRequest):
    determine_seed(request.text_seed_input)
    promot = ""
    if request.customer_promot_flag:
        # 使用自定义文本
        promot = request.customer_promot
    else:
        # 使用文本种子值生成文本
        promot = f'[speed_{request.speed}]'
    params_infer_code = {
        'spk_emb': rand_spk,
        'prompt': promot,
        'temperature': request.temperature,
        'top_P': request.top_P,
        'top_K': request.top_K,
    }
    params_refine_text = {
        'prompt': f'[oral_{request.refine_oral}][laugh_{request.refine_laugh}][break_{request.refine_break}]'}
    return params_infer_code, params_refine_text


@router.post("/generate_audio")
async def generate_audio(request: GenerateAudioRequest):
    chat = create_chat()
    print_current_device_info()
    if request.customer_speaker_flag and request.customer_speaker_tensor:
        # 自定义音色
        speak_tensor = torch.tensor(request.customer_speaker_tensor)
    else:
        # 使用音色种子值创建音色
        speak_tensor = build_random_speaker(chat, request)
    # 生成编码参数
    params_infer_code, params_refine_text = build_code_text(chat, speak_tensor, request)
    # 是否优化文本
    refine_text_flag = request.refine_text_flag
    # 对每个段落分别生成语音
    paragraphs = re.split(r"\n+", request.text)
    text_data_all = []
    audio_data_all = []
    for paragraph in paragraphs:
        if paragraph.replace("\n", "").strip() == '':
            continue
        if refine_text_flag:
            paragraph = chat.infer(paragraph,
                                   skip_refine_text=False,
                                   refine_text_only=True,
                                   params_refine_text=params_refine_text,
                                   params_infer_code=params_infer_code
                                   )
        content = paragraph[0] if isinstance(paragraph, list) else paragraph
        text_data_all.append(content)

        # 生成音频文件
        wav = chat.infer(paragraph,
                         skip_refine_text=True,
                         params_refine_text=params_refine_text,
                         params_infer_code=params_infer_code
                         )
        audio_data = np.array(wav[0]).flatten()
        audio_data_all.append(audio_data)

    # 将所有的语音片段合并成一个完整的音频
    audio_data_all = np.concatenate(audio_data_all)
    sample_rate = 24000
    # 写入WAV文件
    file_name = wav_utils.build_wav_name()
    file_path = str(Path(global_info.AUDIO_TEMP_PATH).joinpath(file_name).resolve())
    wavfile.write(file_path, sample_rate, audio_data_all)

    # 清理显存
    torch.cuda.empty_cache()
    return {
        "out_put_text": "\n".join(text_data_all),
        "audio_name": file_name,
        "tensor": speak_tensor.tolist()
    }


@router.post("/clean/temp/wav")
async def generate_audio():
    wavs = Path(global_info.AUDIO_TEMP_PATH).glob("*.wav")
    # 删除wavs
    for wav in wavs:
        wav.unlink()
    return {
        "code": 0,
        "msg": "ok"
    }

