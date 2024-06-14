<template>
	<div class="system-role-container">
		<el-card shadow="hover">
			<div class="system-user-search mb15">
				<el-form :inline="true" :model="tableData.param">
					<el-form-item label="声音名称">
						<el-input size="default" placeholder="请输入声音名称" style="width: 240px" v-model="tableData.param.name" clearable />
					</el-form-item>
					<el-form-item label="性别">
						<el-select size="default" placeholder="请选择性别" style="width: 240px" v-model="tableData.param.gender" clearable>
							<el-option label="不限" value="0" />
							<el-option label="男" value="1" />
							<el-option label="女" value="2" />
						</el-select>
					</el-form-item>
					<el-form-item>
						<el-button size="default" type="primary" class="ml10" @click="voiceList">
							<el-icon>
								<ele-Search />
							</el-icon>
							查询
						</el-button>
						<el-button size="default" type="success" class="ml10" @click="downloadAudio">
							<el-icon>
								<ele-Download />
							</el-icon>
							下载
						</el-button>
						<el-upload
							ref="upload"
							:on-change="jsonUpload"
							:before-upload="beforeJsonUpload"
							:auto-upload="false"
							:multiple="false"
							:show-file-list="false"
							accept="application/json"
							style="display: flex; align-items: center"
						>
							<template #trigger
								><el-button size="default" type="success" class="ml10" style="flex: 1">
									<el-icon>
										<ele-UploadFilled />
									</el-icon>
									上传
								</el-button></template
							>
						</el-upload>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="tableData.data" style="width: 100%">
				<el-table-column type="index" label="序号" width="60" />
				<el-table-column prop="name" label="音色名称" show-overflow-tooltip></el-table-column>
				<el-table-column prop="gender" label="性别" show-overflow-tooltip>
					<template #default="scope">
						<el-tag v-if="scope.row.gender === 1" type="success">男</el-tag>
						<el-tag v-else-if="scope.row.gender === 2" type="danger">女</el-tag>
						<el-tag v-else type="info">未知</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="describe" label="描述" show-overflow-tooltip></el-table-column>
				<el-table-column prop="voicePath" label="样例音频" show-overflow-tooltip>
					<template #default="scope">
						<audio v-if="scope.row.voicePath !== ''" controls :src="scope.row.voicePath" style="width: 100%"></audio>
					</template>
				</el-table-column>
				<el-table-column label="评分" show-overflow-tooltip>
					<template #default="{ row }">
						<el-rate v-model="row.rate" @change="updateRate(row)" allow-half></el-rate>
					</template>
				</el-table-column>
				<el-table-column prop="count" label="使用次数" show-overflow-tooltip></el-table-column>
				<el-table-column prop="lastAccessTime" label="最后使用时间" show-overflow-tooltip></el-table-column>
				<el-table-column label="操作" width="220">
					<template #default="scope">
						<el-button size="small" text type="primary" @click="genAudio(scope.row)">生成音频</el-button>
						<el-button size="small" text type="primary" @click="downloadAudio(scope.row)">下载</el-button>
						<el-button size="small" text type="primary" @click="onRowDel(scope.row)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
			<pagination
				v-show="tableData.total > 0"
				:total="tableData.total"
				v-model:page="tableData.param.pageNum"
				v-model:limit="tableData.param.pageSize"
				@pagination="voiceList"
			/>
		</el-card>
	</div>
</template>

<script lang="ts">
import { toRefs, reactive, onMounted, ref, defineComponent } from 'vue';
import { ElMessageBox, ElMessage, FormInstance } from 'element-plus';
import { ElNotification } from 'element-plus';
import { getVoiceList, voiceRate, delVoice, voiceGen, voiceDownload, voiceUpload } from '/@/api/cts/voice';
import { ElLoading } from 'element-plus';
import { nextTick } from 'vue';
import type { UploadInstance, UploadProps, UploadRawFile } from 'element-plus';
import { Upload } from '@element-plus/icons-vue/dist/types';
// 定义接口来定义对象的类型
interface TableData {
	id: number;
	name: string;
	gender: string;
	describe: string;
	createdAt: string;
	count: number;
	rate: number;
}
interface TableDataState {
	tableData: {
		data: Array<TableData>;
		total: number;
		param: {
			name: string;
			gender: string;
			pageNum: number;
			pageSize: number;
		};
	};
}

export default defineComponent({
	name: 'apiV1CtsVoiceList',
	setup() {
		const state = reactive<TableDataState>({
			tableData: {
				data: [],
				total: 0,
				param: {
					name: '',
					gender: '0',
					pageNum: 1,
					pageSize: 10,
				},
			},
		});
		//上传
		const uploadUrl = ref('');
		const upload = ref<UploadInstance>();
		// 初始化表格数据
		const initTableData = () => {
			voiceList();
		};
		const voiceList = () => {
			const apiurl = import.meta.env.VITE_API_URL?.endsWith('/') ? import.meta.env.VITE_API_URL.slice(0, -1) : import.meta.env.VITE_API_URL;
			getVoiceList(state.tableData.param).then((res) => {
				state.tableData.data = res.data.voiceList;
				// 遍历state.tableData.data,如果voicePath不为空,则将voicePath
				state.tableData.data.forEach((item: any) => {
					if (item.voicePath) {
						item.voicePath = apiurl + '/api/v1/cts/voice/play?voiceId=' + item.voiceId + '&t=' + new Date().getTime();
					}
				});
				state.tableData.total = res.data.total;
			});
		};
		// 删除声音
		const onRowDel = (row: any) => {
			ElMessageBox.confirm('此操作将永久删除该音色, 是否继续?', '提示', {
				confirmButtonText: '确定',
				cancelButtonText: '取消',
				type: 'warning',
			}).then(() => {
				delVoice(row.voiceId).then((res: any) => {
					if (res.code == 0) {
						ElMessage.success('删除成功');
						initTableData();
					}
				});
			});
		};
		// 评分
		const updateRate = (row: any) => {
			voiceRate(row.voiceId, row.rate).then((res: any) => {
				if (res.code == 0) {
					ElMessage.success('评分修改完成');
				}
			});
		};
		//生成样例音频
		const genAudio = (row: any) => {
			ElMessageBox.confirm('你确定要生成样例音频么?', '提示', {
				confirmButtonText: '确定',
				cancelButtonText: '取消',
				type: 'warning',
			}).then(() => {
				let data = {
					voiceId: row.voiceId,
				};
				const loadingInstance = ElLoading.service({
					lock: true,
					text: '样例音频生成中,请稍后',
					background: 'rgba(0, 0, 0, 0.7)',
				});
				voiceGen(data)
					.then((res: any) => {
						if (res.code === 0) {
							ElNotification({
								title: '音频生成成功',
								message: '音色 ' + row.name + ' 样例音频生成成功',
								type: 'success',
							});
							initTableData();
						}
					})
					.catch((err) => {
						ElNotification({
							title: '音频生成失败',
							message: '音色 ' + row.name + ' 样例音频生成失败',
							type: 'success',
						});
					})
					.finally(() => {
						nextTick(() => {
							loadingInstance.close();
						});
					});
			});
		};
		//下载音频
		const downloadAudio = (row: any) => {
			const apiurl = import.meta.env.VITE_API_URL?.endsWith('/') ? import.meta.env.VITE_API_URL.slice(0, -1) : import.meta.env.VITE_API_URL;
			if (row == null) {
				window.open(apiurl + '/api/v1/cts/voice/download');
			} else {
				window.open(apiurl + '/api/v1/cts/voice/download?voiceId=' + row.voiceId);
			}
		};
		// 上传前准备
		const beforeJsonUpload: UploadProps['beforeUpload'] = (rawFile) => {
			console.log(uploadUrl.value);
			if (rawFile.type !== 'application/json') {
				ElMessage.error('请上传json文件');
				return false;
			}
			return true;
		};
		// 上传
		const jsonUpload = (file: any) => {
			// 创建FormData对象
			const formData = new FormData();
			formData.append('file', file.raw); // file.raw是原始文件对象
			// 发送请求
			voiceUpload(formData).then((res: any) => {
				if (res.code == 0) {
					ElMessage.success('上传成功');
					initTableData();
				}
			});
		};
		// 分页改变
		const onHandleSizeChange = (val: number) => {
			state.tableData.param.pageSize = val;
		};
		// 分页改变
		const onHandleCurrentChange = (val: number) => {
			state.tableData.param.pageNum = val;
		};
		// 页面加载时
		onMounted(() => {
			let baseUrl = import.meta.env.VITE_API_URL?.endsWith('/') ? import.meta.env.VITE_API_URL.slice(0, -1) : import.meta.env.VITE_API_URL;
			uploadUrl.value = baseUrl + '/api/v1/cts/voice/upload';
			initTableData();
		});
		return {
			upload,
			uploadUrl,
			jsonUpload,
			voiceList,
			onRowDel,
			onHandleSizeChange,
			onHandleCurrentChange,
			updateRate,
			genAudio,
			downloadAudio,
			beforeJsonUpload,
			...toRefs(state),
		};
	},
});
</script>
