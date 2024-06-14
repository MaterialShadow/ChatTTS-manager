from datetime import datetime


def build_wav_name():
    return "%s.wav" % datetime.now().strftime("%Y_%m_%d_%H_%M_%S")
