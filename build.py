import logging as log
from pathlib import Path
import shutil
import subprocess
log.basicConfig(level=log.INFO,format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')


def check_npm():
    try:
        # 执行npm --version命令，尝试获取npm版本信息
        result = subprocess.run(['npm', '--version'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True,shell=True)
        # 如果命令执行成功，打印npm版本信息
        log.info(f"npm已安装。版本号：{result.stdout.strip()}")
        return True
    except Exception as e:
        # 如果系统中没有npm命令，打印提示信息
        log.info("系统中未找到npm命令。请确保npm已安装。")
        subprocess.run(['start', 'https://github.com/coreybutler/nvm-windows/releases'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True,shell=True)
        return False

def check_go():
    try:
        # 执行go version命令，尝试获取Go版本信息
        result = subprocess.run(['go', 'version'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True)
        # 如果命令执行成功，打印Go版本信息
        log.info(f"Go已安装。版本信息：{result.stdout.strip()}")
        return True
    except Exception as e:
        subprocess.run(['start', 'https://golang.google.cn/'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True,shell=True)
        log.info(f"发生错误：{e}")
        return False


def buildFrontend(cwd):
    # 构建前端
    log.info("开始构建前端")
    try:
        # 执行npm run build命令
        log.info("开始执行npm构建命令...")
        install_result = subprocess.run(['npm','install'],cwd=cwd,stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True,shell=True)
        if install_result.returncode != 0:
            log.info("npm install 命令执行失败，错误信息如下：")
            log.info(install_result.stderr)
            return

        result = subprocess.run(['npm','run','build'],cwd=cwd,stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True,shell=True)
        # 打印构建过程中的标准输出和标准错误
        log.info("构建完成。以下是输出结果：")
        log.info(result.stdout)
        if result.stderr:
            log.info("以下是错误信息：")
            log.info(result.stderr)
    except subprocess.CalledProcessError as e:
        # 命令执行失败时打印错误信息
        log.info("构建前端过程中发生错误：")
        log.info(e.stderr.strip())
        raise e
    except Exception as e:
        # 其他异常情况
        log.info(f"执行npm命令时发生错误：{e}")
        raise e


def build_go(cwd):
    # 构建Go
    try:
        # 执行go mod tidy命令
        log.info("开始执行go mod tidy...")
        tidy_result = subprocess.run(['go', 'mod', 'tidy'], cwd=cwd,check=True, text=True, capture_output=True)
        log.info("go mod tidy执行完成。")

        # 检查go mod tidy命令的输出，如果有错误则打印并退出
        if tidy_result.returncode != 0:
            log.info("go mod tidy命令执行失败，错误信息如下：")
            log.info(tidy_result.stderr)
            return

        # 执行go build命令
        log.info("开始执行go build...")
        build_result = subprocess.run(['go', 'build','-o','tts.exe'],cwd="managerment", check=True, text=True, capture_output=True,shell=True)
        log.info("go build执行完成。")

        # 打印构建过程中的标准输出和标准错误
        if build_result.stdout:
            log.info("构建输出：")
            log.info(build_result.stdout)
        if build_result.stderr:
            log.info("构建错误：")
            log.info(build_result.stderr)

    except subprocess.CalledProcessError as e:
        # 命令执行失败时打印错误信息
        log.info("构建go过程中发生错误：")
        log.info(e.stderr.strip())
    except Exception as e:
        # 其他异常情况
        log.info(f"执行Go命令时发生错误：{e}")


def copy_base_dir(chatts_dir):
    '''
    复制基础文件
    '''
    manifest_dir = chatts_dir.joinpath("manifest")
    resource_dir = chatts_dir.joinpath("resource")
    hack_dir = chatts_dir.joinpath("hack")
    temp_dir = chatts_dir.joinpath("audio_temp")
    audio_dir = chatts_dir.joinpath("audio")
    db_dir = chatts_dir.joinpath("db")
    if not temp_dir.exists():
        temp_dir.mkdir()
    if not audio_dir.exists():
        audio_dir.mkdir()
    if not db_dir.exists():
        db_dir.mkdir()
    if not manifest_dir.exists():
        log.info("manifest目录不存在，开始拷贝")
        shutil.copytree(Path("managerment").joinpath("manifest"),manifest_dir)
        shutil.move(manifest_dir.joinpath("config").joinpath("config.yaml"),manifest_dir.joinpath("config").joinpath("config.yaml.bak"))
        shutil.move(manifest_dir.joinpath("config").joinpath("config.yaml.final"),manifest_dir.joinpath("config").joinpath("config.yaml"))
    if not resource_dir.exists():
        log.info("resource目录不存在，开始拷贝")
        shutil.copytree(Path("managerment").joinpath("resource"),resource_dir)
        # 复制数据库文件
        shutil.copy(Path("managerment").joinpath("resource").joinpath("data").joinpath("sqlite").joinpath("tts.db"),db_dir)
        log.info("复制数据库文件完成")
    if not hack_dir.exists():
        log.info("hack目录不存在，开始拷贝")
        shutil.copytree(Path("managerment").joinpath("hack"),hack_dir)
    log.info("拷贝基础文件完成")




def copy_front_end(chatts_dir):
    '''
    复制前端资源
    '''
    # 复制前端
    sys_dir = chatts_dir.joinpath("resource").joinpath("public").joinpath("sys")
    # 删除sys_dir内所有文件
    if sys_dir.exists():
        shutil.rmtree(sys_dir)
        log.info("删除sys_dir内所有文件")
    shutil.copytree("managerment/web/dist",sys_dir)
    log.info("拷贝managerment/web/dist完成")


def copy_go_file(chatts_dir):
    '''
    复制打包后的go文件
    '''
    main_exe = chatts_dir.joinpath("tts.exe")
    if main_exe.exists():
        main_exe.unlink()
        log.info("删除tts.exe")
    file_path = chatts_dir.joinpath("managerment").joinpath("tts.exe")
    log.info(file_path.resolve())
    shutil.move(chatts_dir.parent.joinpath("managerment").joinpath("tts.exe"),chatts_dir)
        

def main():
    npm_install_flag = check_npm()
    go_isstall_flag = check_go()
    if not npm_install_flag or not go_isstall_flag:
        log.info("npm或Go未安装，请检查")
        return
    chatts_dir = Path("chattts")
    if not chatts_dir.exists():
        log.info("chatts文件夹不存在,创建chatts文件夹")
        chatts_dir.mkdir()
    else:
        log.info("chatts文件夹已经存在")
    # 复制资源
    copy_base_dir(chatts_dir)
    #构建前端
    managerment_dir = Path("managerment")
    buildFrontend(cwd=managerment_dir.joinpath("web"))
    copy_front_end(chatts_dir)
    # 构建go
    build_go(cwd = managerment_dir)
    copy_go_file(chatts_dir)
    log.info("构建完成")
        

if __name__=="__main__":
    main()