from pathlib import Path

from fastapi import FastAPI, Request

from app.common import global_info
# 路由分组
from app.controller.audio_controller import router as audio_router
import logging as log
log.basicConfig(level=log.INFO,format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')

app = FastAPI()
app.include_router(prefix="/chat", router=audio_router)


def init():
    global_info.AUDIO_TEMP_PATH = str(Path("../chattts/audio_temp").resolve())
    Path(global_info.AUDIO_TEMP_PATH).mkdir(parents=True, exist_ok=True)

if __name__ == "__main__":
    init()

    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8087)
