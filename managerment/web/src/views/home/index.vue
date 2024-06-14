<template>
	<el-card class="github-card">
		<h3>ChatTTS on GitHub</h3>
		<p>Explore the ChatTTS project on GitHub for a generative speech model designed for daily dialogue.</p>
		<el-button type="primary" @click="goToGitHub"> Visit GitHub </el-button>
	</el-card>
	<el-card class="github-card">
		<h3>ChatTTS</h3>
		<p>Free Online ChatTTS</p>
		<el-button type="primary" @click="goToTTS"> Visit chattts </el-button>
	</el-card>
</template>

<script setup>
import { useRouter } from 'vue-router';

const router = useRouter();

function goToGitHub() {
	window.open('https://github.com/2noise/ChatTTS', '_blank');
}
function goToTTS() {
	window.open('https://chattts.com/', '_blank');
}
</script>

<style scoped>
.github-card {
	/* 你的样式 */
}
</style>
