import request from '/@/utils/request';
export function getVoiceList(query: Object) {
	return request({
		url: '/api/v1/cts/voice/list',
		method: 'get',
		params: query,
	});
}

export function getVoiceInfo(voiceId: number) {
	return request({
		url: '/api/v1/cts/voice/get',
		method: 'get',
		params: { voiceId },
	});
}

export function addVoice(data: object) {
	return request({
		url: '/api/v1/cts/voice/add',
		method: 'post',
		data: data,
	});
}

export function editVoice(data: object) {
	return request({
		url: '/api/v1/cts/voice/edit',
		method: 'put',
		data: data,
	});
}

export function delVoice(voiceId: number) {
	return request({
		url: '/api/v1/cts/voice/delete',
		method: 'delete',
		params: { voiceId },
	});
}

export function batchDelVoice(ids: number[]) {
	return request({
		url: '/api/v1/cts/voice/batchDel',
		method: 'delete',
		data: ids,
	});
}

export function voiceRate(voiceId: number, rate: number) {
	return request({
		url: '/api/v1/cts/voice/rate',
		method: 'put',
		data: {
			voiceId: voiceId,
			rate: rate,
		},
	});
}

export function voiceGen(data: object) {
	return request({
		url: '/api/v1/cts/voice/gen',
		method: 'post',
		data: data,
	});
}

export function voiceDownload(voiceId: number) {
	return request({
		url: '/api/v1/cts/voice/download',
		method: 'get',
		params: { voiceId },
	});
}

export function voiceUpload(data: object) {
	return request({
		url: '/api/v1/cts/voice/upload',
		method: 'post',
		data: data,
		headers: {
			'Content-Type': 'multipart/form-data',
		},
	});
}
