<template>
	<div class="profile-page">
		<el-card class="profile-card" shadow="never">
			<div class="profile-header">
				<el-avatar :size="88" :src="profileImage">
					{{ avatarText }}
				</el-avatar>
				<div class="profile-identity">
					<h1>个人中心</h1>
					<p>{{ displayName }} · {{ displayEmail }}</p>
				</div>
			</div>

			<el-divider />

			<div class="profile-nav-wrap">
				<el-menu
					:default-active="activeMenu"
					mode="horizontal"
					class="profile-menu"
					@select="handleSelect"
				>
					<el-menu-item index="profile-info">个人信息</el-menu-item>
					<el-menu-item index="profile-posts">我的帖子</el-menu-item>
					<el-menu-item index="profile-relations">关系</el-menu-item>
				</el-menu>
			</div>

			<div class="profile-content">
				<router-view />
			</div>
		</el-card>
	</div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useAuthStore } from '../stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const user = computed(() => authStore.user || {})
const displayName = computed(() => user.value?.name || '未命名用户')
const displayEmail = computed(() => user.value?.email || '暂无邮箱信息')
const profileImage = computed(() => user.value?.imageUrl || '')

const activeMenu = computed(() => {
	if (route.name === 'profile-posts') {
		return 'profile-posts'
	}
	if (route.name === 'profile-relations') {
		return 'profile-relations'
	}
	return 'profile-info'
})

const avatarText = computed(() => {
	const name = displayName.value.trim()
	return name ? name.slice(0, 1).toUpperCase() : 'U'
})

function handleSelect(menuName) {
	router.push({ name: menuName })
}

onMounted(() => {
	authStore.hydrate()
})
</script>

<style scoped>
.profile-page {
	min-height: calc(100vh - 72px);
	padding: 28px 20px;
	background:
		radial-gradient(circle at top right, rgba(37, 99, 235, 0.12), transparent 36%),
		radial-gradient(circle at bottom left, rgba(16, 185, 129, 0.12), transparent 32%),
		linear-gradient(180deg, #f7fbff 0%, #eef6ff 100%);
}

.profile-card {
	max-width: 920px;
	margin: 0 auto;
	padding: 6px;
	border-radius: 24px;
	border: 1px solid rgba(148, 163, 184, 0.14);
	box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08);
}

.profile-header {
	display: flex;
	align-items: center;
	gap: 18px;
	padding: 8px 6px;
}

.profile-identity h1 {
	margin: 0;
	font-size: clamp(1.6rem, 3vw, 2.1rem);
	line-height: 1.15;
	color: #0f172a;
}

.profile-identity p {
	margin: 8px 0 0;
	color: #64748b;
}

.profile-nav-wrap {
	padding: 0 8px;
}

.profile-menu {
	border-bottom: 0;
	background: transparent;
}

.profile-menu :deep(.el-menu-item) {
	font-weight: 600;
	height: 44px;
	line-height: 44px;
}

.profile-content {
	margin-top: 12px;
}

@media (max-width: 760px) {
	.profile-page {
		padding: 16px 12px;
	}

	.profile-header {
		flex-direction: column;
		text-align: center;
	}
}
</style>