<template>
	<div class="splitpanes-container">
		<el-card shadow="hover" header="ChatTTS WEBUI">
			<splitpanes class="default-theme" @resize="paneSize = $event[0].size">
				<pane :size="paneSize">
					<el-tabs type="border-card">
						<el-tab-pane label="文本输入">
							<el-form :model="form" label-position="top" label-width="90%">
								<el-row :gutter="20">
									<el-col :span="22">
										<el-input v-model="textContent" :rows="textAreaRows" type="textarea" placeholder="请输入文本内容"> </el-input>
									</el-col>
									<el-col :span="2">
										<el-button large class="gen-audio-button" type="primary" @click="genAudio">音频生成</el-button>
									</el-col>
								</el-row>
								<el-row :gutter="20" style="margin-top: 20px">
									<el-col :span="8">
										<el-button type="primary" :icon="Upload" @click="openFileSelector">文本上传</el-button>
										<el-button type="primary" @click="removeWhitespaceLines">去除空行</el-button>
										<el-button type="primary" @click="removeWhitespaceLines">待施工</el-button>
									</el-col>
									<!-- 隐藏的文件输入元素，用于触发文件选择 -->
									<input type="file" ref="fileInput" style="display: none" @change="handleFileChange" />
								</el-row>
							</el-form>
						</el-tab-pane>
						<el-tab-pane label="文本输出">
							<el-form :model="form" label-position="top" label-width="90%">
								<el-row :gutter="20">
									<el-col :span="24">
										<el-input v-model="outputText" :rows="textAreaRows" type="textarea" placeholder="请输入文本内容"> </el-input>
									</el-col>
								</el-row>
							</el-form>
						</el-tab-pane>
						<el-tab-pane label="使用说明">
							<el-form :model="form" label-position="top" label-width="90%">
								<el-row :gutter="20">
									<el-col :span="24">
										<el-input v-model="illustrate" :rows="textAreaRows" type="textarea" placeholder="请输入文本内容" disabled> </el-input>
									</el-col>
								</el-row>
								<!-- <el-row :gutter="20">
									<el-col :span="8">
										<el-button :icon="Download" type="primary">文本下载</el-button>
									</el-col>
								</el-row> -->
							</el-form>
						</el-tab-pane>
					</el-tabs>

					<!--音频输出-->
					<el-card style="width: 100%">
						<template #header>
							<div class="card-header">
								<span>
									<el-text class="mx-1" type="info">音频输出</el-text>
									&nbsp;&nbsp;
									<el-text class="mx-1" type="info">{{ durationText }}</el-text>
								</span>
							</div>
						</template>
						<audio controls :src="audioSrc" style="width: 100%"></audio>
					</el-card>

					<!--快捷操作-->
					<el-card style="width: 100%">
						<template #header>
							<div class="card-header">
								<span>
									<el-text class="mx-1" type="info">快捷操作</el-text>
								</span>
							</div>
						</template>
						<el-form :model="form">
							<el-row :gutter="20">
								<el-col :span="2">
									<el-button type="primary" @click="cleanTempWavs">清理缓存</el-button>
								</el-col>
								<el-col :span="2">
									<el-button type="primary" @click="openDialog" :disabled="currentTensor.length === 0">保存音色</el-button>
								</el-col>
							</el-row>
						</el-form>
					</el-card>
					<AddVoiceDialog v-model="voiceDialogVisible" @submit="handleDialogSubmit" />
				</pane>
				<pane>
					<!-- 输入区域 -->
					<el-form :model="form" label-position="top" label-width="90%">
						<!--核心音频控制-->
						<el-card style="width: 100%">
							<template #header>
								<div class="card-header">
									<span><el-text class="mx-1" type="info">核心音频控制</el-text></span>
								</div>
							</template>
							<el-row :gutter="20">
								<el-col :span="8">
									<el-switch v-model="form.refine_text_flag" class="mb-2" active-text="开启文本优化" inactive-text="关闭文本优化" inline-prompt />
								</el-col>
								<el-col :span="8">
									<el-switch
										v-model="form.customer_promot_flag"
										class="mb-2"
										active-text="启用自定义prmot"
										inactive-text="关闭自定义promot"
										inline-prompt
									></el-switch>
								</el-col>
								<el-col :span="8">
									<el-switch
										v-model="form.customer_speaker_flag"
										class="mb-2"
										active-text="启用自定义音色"
										inactive-text="关闭自定义音色"
										inline-prompt
									/>
								</el-col>
								<el-col :span="8" style="margin-top: 20px">
									<el-select
										v-model="form.customer_speaker_id"
										filterable
										:filter-method="queryVoiceList"
										placeholder="请选择音色"
										style="width: 240px"
										clearable
										v-if="form.customer_speaker_flag"
									>
										<el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
									</el-select>
								</el-col>
							</el-row>
						</el-card>
						<!--口语化控制-->
						<el-card style="width: 100%" v-show="form.refine_text_flag">
							<template #header>
								<div class="card-header">
									<span><el-text class="mx-1" type="info">口语化控制</el-text></span>
								</div>
							</template>
							<el-row :gutter="20">
								<el-col :span="8">
									<el-form-item>
										<el-tooltip content="控制文本口语化程度，i范围为 0-9，数字越大，添加的“啊”、“就是”、“那”之类的口头词越多,默认值2" placement="top">
											<el-tag type="info">口语化(oral)</el-tag>
										</el-tooltip>
										<el-slider v-model="form.refine_oral" :min="0" :max="9" :step="1" show-input></el-slider>
									</el-form-item>
								</el-col>
								<el-col :span="8">
									<el-form-item>
										<el-tooltip content="控制文本中添加笑声的程度，i范围为 0-9，值越大，笑声越多，默认值0" placement="top">
											<el-tag type="info">笑(laugh)</el-tag>
										</el-tooltip>
										<el-slider v-model="form.refine_laugh" :min="0" :max="9" :step="1" show-input></el-slider>
									</el-form-item>
								</el-col>
								<el-col :span="8">
									<el-form-item>
										<el-tooltip content="控制文本中添加停顿的程度，i范围为 0-9，值越大，停顿越多" placement="top">
											<el-tag type="info">停顿(break)</el-tag>
										</el-tooltip>
										<el-slider v-model="form.refine_break" :min="0" :max="9" :step="1" show-input></el-slider>
									</el-form-item>
								</el-col>
							</el-row>
						</el-card>
						<!--自定义promot-->
						<el-card style="width: 100%" v-show="form.customer_promot_flag">
							<template #header>
								<div class="card-header">
									<span><el-text class="mx-1" type="info">自定义promot</el-text></span>
								</div>
							</template>
							<el-row :gutter="20">
								<el-col :span="8">
									<el-form-item>
										<el-input v-model="form.customer_prompt"></el-input>
									</el-form-item>
								</el-col>
							</el-row>
						</el-card>
						<!--速度控制-->
						<el-card style="width: 100%">
							<template #header>
								<div class="card-header">
									<span><el-text class="mx-1" type="info">速度</el-text></span>
								</div>
							</template>
							<el-row :gutter="20">
								<el-col :span="12">
									<el-form-item>
										<el-tooltip content="速度" placement="top">
											<el-tag type="info">速度(speed)</el-tag>
										</el-tooltip>
										<el-slider v-model="form.speed" :min="0" :max="9" :step="1" show-input></el-slider>
									</el-form-item>
								</el-col>
							</el-row>
						</el-card>

						<!--音频采样控制-->
						<el-card style="width: 100%" v-if="!form.customer_speaker_flag">
							<template #header>
								<div class="card-header">
									<span><el-text class="mx-1" type="info">随机种子</el-text></span>
								</div>
							</template>
							<el-row :gutter="20">
								<el-col :span="12">
									<el-form-item>
										<el-tooltip content="随机音频种子" placement="top">
											<el-tag type="info">随机音频种子(audio_seed)</el-tag>
										</el-tooltip>
										<el-input v-model="form.audio_seed_input" type="number" placeholder="请输入音频种子" class="input-with-select">
											<template #append>
												<el-button type="primary" @click="generateAudioRandomNumber">随机</el-button>
											</template>
										</el-input>
									</el-form-item>
								</el-col>

								<el-col :span="12">
									<el-form-item>
										<el-tooltip content="随机文本种子" placement="top">
											<el-tag type="info">随机文本种子(text_seed)</el-tag>
										</el-tooltip>
										<el-input v-model="form.text_seed_input" type="number" placeholder="请输入文本种子" class="input-with-select">
											<template #append>
												<el-button type="primary" @click="generateTextRandomNumber">随机</el-button>
											</template>
										</el-input>
									</el-form-item>
								</el-col>
							</el-row>
						</el-card>
						<!--音频采样控制-->
						<el-card style="width: 100%">
							<template #header>
								<div class="card-header">
									<span><el-text class="mx-1" type="info">音频采样控制</el-text></span>
								</div>
							</template>
							<el-row :gutter="20">
								<el-col :span="8">
									<el-form-item>
										<el-tooltip
											content="这个参数控制着生成文本时的随机性。temperature的值越低，生成的文本就越倾向于选择概率最高的词汇，从而更加确定性；值越高，生成的文本就越随机"
											placement="top"
										>
											<el-tag type="info">音频采样温度(temperature)</el-tag>
										</el-tooltip>
										<el-slider v-model="form.temperature" :min="0.00001" :max="1.0" :step="0.00001" show-input></el-slider>
									</el-form-item>
								</el-col>

								<el-col :span="8">
									<el-form-item>
										<el-tooltip
											content="top_P是一个概率阈值，意味着模型将只考虑累积概率达到这个阈值的词汇。例如，如果设置为0.7，那么模型将只考虑前30%的词汇。调大这个值有助于生成更多样化的文本,默认0.7"
											placement="top"
										>
											<el-tag type="info">音频采样概率阈值(top_P)</el-tag>
										</el-tooltip>
										<el-slider v-model="form.top_P" :min="0.1" :max="0.9" :step="0.05" show-input></el-slider>
									</el-form-item>
								</el-col>

								<el-col :span="8">
									<el-form-item>
										<el-tooltip
											content="与top_P类似，但它是直接限制生成词汇的数量。top_K表示模型在生成下一个词时只考虑概率最高的K个词汇"
											placement="top"
										>
											<el-tag type="info">音频采样词汇量(top_K)</el-tag>
										</el-tooltip>
										<el-slider v-model="form.top_K" :min="1" :max="20" :step="0.1" show-input></el-slider>
									</el-form-item>
								</el-col>
							</el-row>
						</el-card>
					</el-form>
				</pane>
			</splitpanes>
		</el-card>
	</div>
</template>

<script lang="ts">
import { ref } from 'vue';
import { defineComponent } from 'vue';
import { Delete, Edit, Search, Share, Upload, Download } from '@element-plus/icons-vue';
import { Splitpanes, Pane } from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';
import 'element-plus/dist/index.css';
import { onMounted } from 'vue';
import { getVoiceList, addVoice } from '/@/api/cts/voice';
import { GenAdudio, CleanWavs } from '/@/api/cts/tts/index';
import { ElLoading } from 'element-plus';
import { nextTick } from 'vue';
import AddVoiceDialog from './component/addVoiceDialog.vue';
import { ElMessage } from 'element-plus';

export default defineComponent({
	name: 'ChatTTSWebUI',
	components: {
		Splitpanes,
		Pane,
		AddVoiceDialog,
	},
	setup() {
		const paneSize = ref(60); // 默认分割比例
		const maxRandomNum = ref(10000);
		// 文本域行数
		const textAreaRows = ref(15);
		const textContent = ref<string>('曾经有一段真诚的爱情摆在我面前,我没有珍惜,直到失去了才后悔莫及');
		const fileInput = ref<HTMLInputElement | null>(null);
		const illustrate = ref(
			'待施工。\n名词解释参考:https://zhuanlan.zhihu.com/p/703240560\n默认音色参考:https://github.com/2noise/ChatTTS/issues/238'
		);

		// 确保audioSrc被定义并具有初始值
		const audioName = ref('');
		const apiUrl = import.meta.env.VITE_API_URL?.endsWith('/') ? import.meta.env.VITE_API_URL.slice(0, -1) : import.meta.env.VITE_API_URL;
		const audioSrc = ref('');
		const durationText = ref('');

		//loading
		const loading = ref(false);
		// 新增音色对话框
		const voiceDialogVisible = ref(false);
		const form = ref({
			text: '',
			speed: 2,
			temperature: 0.3,
			top_P: 0.7,
			top_K: 20,
			refine_oral: 2,
			refine_laugh: 0,
			refine_break: 0,
			audio_seed_input: 2,
			text_seed_input: 42,
			refine_text_flag: false,
			customer_speaker_flag: false,
			customer_speaker_id: null,
			customer_promot_flag: false,
			customer_prompt: '[speed_2]',
			// 其他表单项的初始值
		});
		// 音色下拉框
		const options = ref([]);
		// 文本输出
		const outputText = ref('');
		// 当前tensor
		const currentTensor = ref([]);
		// 定义一个函数来生文本数字
		const generateAudioRandomNumber = () => {
			// 可以根据需要生成不同范围的随机数，这里以1到10000为例
			const randomNumber = Math.floor(Math.random() * maxRandomNum.value) + 1;
			form.value.audio_seed_input = randomNumber; // 设置输入框的值为随机数
		};
		const generateTextRandomNumber = () => {
			// 可以根据需要生成不同范围的随机数，这里以1到10000为例
			const randomNumber = Math.floor(Math.random() * maxRandomNum.value) + 1;
			form.value.text_seed_input = randomNumber; // 设置输入框的值为随机数
		};
		const genAudio = () => {
			if (loading.value == true) {
				return;
			}
			audioName.value = '';
			audioSrc.value = '';
			durationText.value = '';
			// 这里应该添加表单提交的逻辑
			console.log('提交');
			form.value.text = textContent.value;
			loading.value = true;
			const loadingInstance = ElLoading.service({
				lock: true,
				text: '音频生成中',
				background: 'rgba(0, 0, 0, 0.7)',
			});
			let startTime = new Date().getTime();
			GenAdudio(form.value)
				.then((res: any) => {
					// 处理响应数据
					let endTime = new Date().getTime(); // 请求完成时间
					let duration = endTime - startTime;
					console.log('耗时' + duration + 'ms');
					durationText.value = '耗时' + duration / 1000 + 's';
					audioName.value = res.data.audio_name.replace('.wav', '');
					let url = apiUrl + '/api/v1/cts/tts/download?audio_name=' + audioName.value;
					audioSrc.value = url;
					outputText.value = res.data.out_put_text;
					currentTensor.value = res.data.tensor;
					return;
				})
				.catch((err: any) => {
					let endTime = new Date().getTime();
					let duration = endTime - startTime;
					durationText.value = '耗时' + duration / 1000 + 's,生成失败';
					return;
				})
				.finally(() => {
					loading.value = false;
					nextTick(() => {
						loadingInstance.close();
					});
				});
		};

		const openFileSelector = () => {
			// 确保ref已经定义并且可以访问
			if (fileInput.value) {
				fileInput.value.click();
			}
		};

		const handleFileChange = (event: Event) => {
			const inputElement = event.target as HTMLInputElement;
			if (inputElement.files && inputElement.files.length > 0) {
				const file = inputElement.files[0];
				const reader = new FileReader();
				reader.onload = (e: ProgressEvent<FileReader>) => {
					textContent.value = e.target?.result as string;
				};
				reader.readAsText(file);
			}
		};

		const queryVoiceList = (name: string) => {
			let data = {
				orderBy: 'last_access_time desc',
				name: name,
			};
			getVoiceList(data).then((res: any) => {
				if (res.code === 0) {
					let optionList: any = [];
					res.data.voiceList.forEach((item: any) => {
						optionList.push({
							value: item.voiceId,
							label: item.name,
						});
					});
					//form.value.customer_speaker_id = optionList[0].value;
					options.value = optionList;
				}
			});
		};

		const openDialog = () => {
			console.log('打开弹窗');
			voiceDialogVisible.value = true;
			console.log(voiceDialogVisible.value);
		};

		const handleDialogSubmit = (formData: any) => {
			// 处理表单提交逻辑
			voiceDialogVisible.value = false;
			let data: any = {};
			Object.assign(data, formData);
			data['tensor'] = currentTensor.value;
			data['audioName'] = audioName.value;
			console.log(data);
			addVoice(data).then((res: any) => {
				if (res.code === 0) {
					ElMessage.success('添加成功');
				}
			});
		};

		const cleanTempWavs = () => {
			CleanWavs().then((res: any) => {
				if (res.code == 0) {
					ElMessage.success('清理成功');
				}
			});
		};

		function removeWhitespaceLines() {
			//去除文本框中的空行
			// 将文本分割成行
			let lines = textContent.value.split('\n');
			// 使用filter方法去除空行
			textContent.value = lines.filter((line) => !/^\s*$/.test(line)).join('\n');
		}

		// 由于使用了setup()，需要使用onMounted来挂载ref
		onMounted(() => {
			fileInput.value = document.querySelector('input[type="file"]');
			queryVoiceList();
		});
		return {
			paneSize,
			form,
			outputText,
			audioSrc,
			textAreaRows,
			textContent,
			Upload,
			Download,
			fileInput,
			options,
			apiUrl,
			loading,
			genAudio,
			audioName,
			durationText,
			cleanTempWavs,
			queryVoiceList,
			handleFileChange,
			openFileSelector,
			currentTensor,
			openDialog,
			handleDialogSubmit,
			illustrate,
			voiceDialogVisible,
			removeWhitespaceLines,
			generateAudioRandomNumber,
			generateTextRandomNumber,
		};
	},
});
</script>

<style scoped lang="scss">
.splitpanes-container {
	width: 100%;
	height: 100vh;
	margin: 0;
	padding: 0;
}

.output-pane {
	padding: 20px;
	height: 100%;
	box-sizing: border-box;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
}

.gen-audio-button {
	width: 100% !important;
	height: 100% !important;
}
</style>
