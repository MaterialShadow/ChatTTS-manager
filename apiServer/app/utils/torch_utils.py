import torch
import logging

def print_current_device_info():
    logger = logging.getLogger(__name__)
    if torch.cuda.is_available():
        available_gpus = []
        for i in range(torch.cuda.device_count()):
            props = torch.cuda.get_device_properties(i)
            free_memory = props.total_memory - torch.cuda.memory_reserved(i)
            available_gpus.append((i, free_memory))
            selected_gpu, max_free_memory = max(available_gpus, key=lambda x: x[1])
            free_memory_mb = max_free_memory / (1024 * 1024)
            logger.log(logging.WARNING, f'GPU {selected_gpu} has {round(free_memory_mb, 2)} MB memory left.')
    else:
        logger.log(logging.WARNING, f'No GPU found')