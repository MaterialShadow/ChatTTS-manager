import request from '/@/utils/longTimeOutRequest';

export function GenAdudio(data: object) {
	return request({
		url: '/api/v1/cts/tts/gen',
		method: 'post',
		data: data,
	});
}

export function CleanWavs() {
	return request({
		url: '/api/v1/cts/tts/clean',
		method: 'DELETE',
	});
}
