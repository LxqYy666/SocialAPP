<template>
	<div class="home-page">
		<el-card class="home-card" shadow="never">
			<div class="home-header">
				<div>
					<p class="eyebrow">Dashboard</p>
					<h2>欢迎回来</h2>
					<p class="description">
						你当前已经登录，可以继续浏览动态、消息和通知。
					</p>
				</div>
				<el-button type="danger" plain @click="handleLogout">退出登录</el-button>
			</div>

			<el-divider />

			<div class="profile-block">
				<el-avatar :size="64">{{ avatarText }}</el-avatar>
				<div>
					<div class="profile-name">{{ displayName }}</div>
					<div class="profile-meta">{{ displayEmail }}</div>
				</div>
			</div>
		</el-card>
	</div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const displayName = computed(() => authStore.user?.name || '未命名用户')
const displayEmail = computed(() => authStore.user?.email || '暂无邮箱信息')
const avatarText = computed(() => {
	const name = displayName.value.trim()
	return name ? name.slice(0, 1).toUpperCase() : 'U'
})

function handleLogout() {
	authStore.logout()
	router.replace({ name: 'auth', params: { mode: 'login' } })
}

onMounted(() => {
	authStore.hydrate()
})
</script>

<style scoped>
.home-page {
	min-height: calc(100vh - 72px);
	padding: 28px 20px;
	background: linear-gradient(180deg, #f8fbff 0%, #eef4ff 100%);
}

.home-card {
	max-width: 920px;
	margin: 0 auto;
	border-radius: 24px;
	border: 1px solid rgba(148, 163, 184, 0.14);
	box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08);
}

.home-header {
	display: flex;
	align-items: flex-start;
	justify-content: space-between;
	gap: 16px;
}

.eyebrow {
	margin: 0 0 8px;
	font-size: 12px;
	font-weight: 700;
	letter-spacing: 0.16em;
	text-transform: uppercase;
	color: #64748b;
}

.home-header h2 {
	margin: 0;
	font-size: 32px;
	line-height: 1.1;
	color: #0f172a;
}

.description {
	margin: 10px 0 0;
	color: #475569;
	line-height: 1.7;
}

.profile-block {
	display: flex;
	align-items: center;
	gap: 16px;
}

.profile-name {
	font-size: 18px;
	font-weight: 700;
	color: #0f172a;
}

.profile-meta {
	margin-top: 4px;
	color: #64748b;
}

@media (max-width: 640px) {
	.home-page {
		padding: 16px 12px;
	}

	.home-header {
		flex-direction: column;
	}

	.home-header h2 {
		font-size: 26px;
	}
}
</style>