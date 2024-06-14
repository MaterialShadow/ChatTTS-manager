<template>
	<el-dialog :title="title" v-model="visible" :before-close="handleClose" width="580px">
		<el-form :model="form" :rules="rules" ref="formRef">
			<el-row :gutter="20">
				<el-col :span="24">
					<el-form-item label="名称" prop="name">
						<el-input v-model="form.name"></el-input>
					</el-form-item>
				</el-col>
				<el-col :span="24">
					<el-form-item label="性别" prop="gender">
						<el-radio-group v-model="form.gender">
							<el-radio :label="1">男</el-radio>
							<el-radio :label="2">女</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :span="24">
					<el-form-item label="描述" prop="describe">
						<el-input type="textarea" v-model="form.describe" row="35"> </el-input>
					</el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<template #footer>
			<el-button @click="handleClose">取消</el-button>
			<el-button type="primary" @click="submitForm">确定</el-button>
		</template>
	</el-dialog>
</template>

<script>
import { ref, defineComponent, watch, onMounted } from 'vue';

export default defineComponent({
	name: 'AddVoiceDialog',
	props: {
		title: {
			type: String,
			default: '新增音色',
		},
		modelValue: {
			type: Boolean,
			default: false,
		},
	},
	emits: ['update:modelValue', 'submit'],
	setup(props, { emit }) {
		const formRef = ref(null); // 表单引用
		const form = ref({
			name: '',
			gender: 1,
			describe: '',
		});
		const visible = ref(props.modelValue);
		const rules = ref({
			name: [
				{ required: true, message: '名称不能为空', trigger: 'blur' },
				{ min: 3, max: 10, message: '长度在 3 到 10 个字符', trigger: 'blur' },
			],
			gender: [{ required: true, message: '性别不能为空', trigger: 'blur' }],
			describe: [{ required: true, message: '描述不能为空', trigger: 'blur' }],
		});

		const handleClose = () => {
			visible.value = false;
			emit('update:modelValue', visible.value);
			formRef.value.resetFields(); // 重置表单
		};

		const submitForm = () => {
			formRef.value.validate((valid) => {
				if (valid) {
					console.log('Form submitted', form.value);
					emit('submit', form.value);
					visible.value = false;
					formRef.value.resetFields();
				} else {
					console.error('表单校验失败');
					return false;
				}
			});
		};

		watch(
			() => props.modelValue,
			(newVal) => {
				visible.value = newVal;
			}
		);

		onMounted(() => {
			// 可以在组件挂载后执行一些初始化操作
		});

		return {
			form,
			visible,
			rules,
			formRef,
			handleClose,
			submitForm,
		};
	},
});
</script>
