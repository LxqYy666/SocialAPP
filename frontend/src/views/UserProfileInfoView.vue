<template>
	<div class="user-profile-info-page">
		<el-skeleton :loading="loading" animated>
			<template #default>
				<template v-if="profileUser">
					<div class="profile-grid">
						<div class="profile-item detail-list">
							<div class="detail-row">
								<span class="label">用户名称</span>
								<span class="detail-value">{{ profileName }}</span>
							</div>
							<div class="detail-row">
								<span class="label">邮箱地址</span>
								<span class="detail-value">{{ profileEmail }}</span>
							</div>
							<div class="detail-row">
								<span class="label">个人简介</span>
								<p class="detail-bio">{{ profileBio }}</p>
							</div>
						</div>
						<div class="profile-item stats">
							<div>
								<span class="value">{{ followersCount }}</span>
								<span class="label">粉丝</span>
							</div>
							<div>
								<span class="value">{{ followingCount }}</span>
								<span class="label">关注</span>
							</div>
						</div>
					</div>

					<div v-if="!isOwnProfile" class="follow-bar">
						<el-button
							type="primary"
							:plain="isFollowing"
							:loading="followLoading"
							@click="toggleFollow"
						>
							{{ isFollowing ? '取消关注' : '关注' }}
						</el-button>
						<el-button type="success" plain @click="openChat">聊天</el-button>
					</div>
				</template>
				<el-empty v-else description="未找到用户信息" :image-size="100" />
			</template>
		</el-skeleton>
	</div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'

import { userApi } from '../api'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const followLoading = ref(false)
const profileUser = ref(null)

const profileId = computed(() => String(route.params.id || ''))
const profileName = computed(() => profileUser.value?.name || '未命名用户')
const profileEmail = computed(() => profileUser.value?.email || '暂无邮箱信息')
const profileBio = computed(() => profileUser.value?.bio || '这个人很神秘，还没有填写简介。')
const followersCount = computed(() => profileUser.value?.followers?.length || 0)
const followingCount = computed(() => profileUser.value?.following?.length || 0)
const currentUserId = computed(() => authStore.userId || '')
const isOwnProfile = computed(() => Boolean(profileId.value && profileId.value === currentUserId.value))
const isFollowing = computed(() => Boolean(authStore.user?.following?.includes(profileId.value)))

async function loadProfile() {
	if (!profileId.value) {
		profileUser.value = null
		return
	}

	loading.value = true
	try {
		const response = await userApi.getById(profileId.value)
		profileUser.value = response?.user || null
	} catch {
		profileUser.value = null
	} finally {
		loading.value = false
	}
}

async function toggleFollow() {
	if (!profileId.value || isOwnProfile.value) {
		return
	}

	followLoading.value = true
	try {
		const response = await userApi.follow(profileId.value)
		if (response?.firstUser) {
			authStore.updateUser(response.firstUser)
		}
		if (response?.secondUser) {
			profileUser.value = response.secondUser
		}
		await loadProfile()
	} catch (error) {
		ElMessage.error(error?.response?.data?.error || '操作失败，请稍后重试')
	} finally {
		followLoading.value = false
	}
}

function openChat() {
	if (!profileId.value || isOwnProfile.value) {
		return
	}

	router.push({ name: 'messages', query: { target: profileId.value } })
}

watch(profileId, loadProfile, { immediate: true })

onMounted(() => {
	authStore.hydrate()
})
</script>

<style scoped>
.user-profile-info-page {
	min-height: 100%;
}

.profile-grid {
	display: grid;
	grid-template-columns: 1.2fr 1fr;
	gap: 16px;
}

.profile-item {
	padding: 18px;
	border-radius: 16px;
	background: #f8fafc;
	border: 1px solid rgba(148, 163, 184, 0.2);
}

.detail-list {
	display: flex;
	flex-direction: column;
	gap: 14px;
}

.detail-row {
	display: flex;
	flex-direction: column;
	gap: 6px;
}

.detail-value {
	font-size: 15px;
	font-weight: 600;
	color: #0f172a;
	word-break: break-all;
}

.detail-bio {
	margin: 0;
	color: #334155;
	line-height: 1.8;
}

.label {
	display: block;
	font-size: 12px;
	font-weight: 700;
	letter-spacing: 0.08em;
	text-transform: uppercase;
	color: #64748b;
}

.stats {
	display: grid;
	grid-template-columns: repeat(2, 1fr);
	gap: 12px;
}

.stats > div {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	min-height: 120px;
	background: #ffffff;
	border-radius: 14px;
	border: 1px solid rgba(148, 163, 184, 0.16);
}

.value {
	font-size: 28px;
	font-weight: 800;
	line-height: 1;
	color: #0f172a;
	margin-bottom: 8px;
}

.follow-bar {
	margin-top: 16px;
	display: flex;
	justify-content: flex-end;
}

@media (max-width: 760px) {
	.profile-grid {
		grid-template-columns: 1fr;
	}
}
</style>