import logging as log
from pathlib import Path
import shutil
log.basicConfig(level=log.INFO,format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')



if __name__=="__main__":
    """
    chatts 目录层级过多,直接删除可能会报需要管理员权限,可以尝试用python删除
    """
    log.info("Starting Delete Chats")
    chattts_dir = Path("chattts")
    shutil.rmtree(chattts_dir)