<template>
	<div class="profile-relations-page">
		<div class="profile-item relations-panel">
			<div class="relations-head">
				<span class="label">关系</span>
				<el-tag size="small" type="info">总计 {{ followersCount + followingCount }}</el-tag>
			</div>

			<el-tabs v-model="relationTab" class="relation-tabs">
				<el-tab-pane :label="`关注中 (${followingUsers.length})`" name="following" />
				<el-tab-pane :label="`粉丝 (${followersUsers.length})`" name="followers" />
			</el-tabs>

			<el-skeleton :loading="relationLoading" :rows="3" animated>
				<template #default>
					<div v-if="activeRelationUsers.length" class="relation-list">
						<div
							v-for="person in activeRelationUsers"
							:key="person.id"
							class="relation-user-card"
							role="button"
							tabindex="0"
							@click="openUserProfile(person.id)"
							@keydown.enter.prevent="openUserProfile(person.id)"
							@keydown.space.prevent="openUserProfile(person.id)"
						>
							<el-avatar :size="42" :src="person.imageUrl">
								{{ person.name.slice(0, 1).toUpperCase() }}
							</el-avatar>
							<div class="relation-user-meta">
								<div class="relation-user-name">{{ person.name }}</div>
								<div class="relation-user-email">{{ person.email || '暂无邮箱' }}</div>
							</div>
							<div class="relation-actions">
								<el-button size="small" plain @click.stop="openUserProfile(person.id)">主页</el-button>
								<el-button size="small" type="primary" @click.stop="openChat(person.id)">聊天</el-button>
							</div>
						</div>
					</div>
					<el-empty v-else description="暂无关系数据" :image-size="90" />
				</template>
			</el-skeleton>
		</div>
	</div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { userApi } from '../api'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const relationLoading = ref(false)
const relationTab = ref('following')
const followingUsers = ref([])
const followersUsers = ref([])

const user = computed(() => authStore.user || {})
const followersCount = computed(() => user.value?.followers?.length || 0)
const followingCount = computed(() => user.value?.following?.length || 0)
const activeRelationUsers = computed(() => (relationTab.value === 'following' ? followingUsers.value : followersUsers.value))

function openUserProfile(userId) {
	if (!userId) {
		return
	}
	router.push({ name: 'user-profile', params: { id: userId } })
}

function openChat(userId) {
	if (!userId) {
		return
	}
	router.push({ name: 'messages', query: { target: userId } })
}

function normalizeRelationUser(rawUser, fallbackId) {
	return {
		id: rawUser?.id || rawUser?._id || fallbackId,
		name: rawUser?.name || '未命名用户',
		email: rawUser?.email || '',
		imageUrl: rawUser?.imageUrl || '',
	}
}

async function fetchUsersByIds(ids = []) {
	if (!ids.length) {
		return []
	}

	const requests = ids.map((id) => userApi.getById(id))
	const result = await Promise.allSettled(requests)

	return result
		.filter((item) => item.status === 'fulfilled' && item.value?.user)
		.map((item) => normalizeRelationUser(item.value.user, item.value.user?.id || item.value.user?._id || ''))
}

async function refreshUser() {
	if (!authStore.userId) {
		return
	}

	try {
		const response = await userApi.getById(authStore.userId)
		if (response?.user) {
			authStore.updateUser(response.user)
		}
	} catch {
		// Ignore refresh error to avoid blocking relation rendering with cached user data.
	}
}

async function refreshRelations() {
	relationLoading.value = true
	try {
		followingUsers.value = await fetchUsersByIds(user.value?.following || [])
		followersUsers.value = await fetchUsersByIds(user.value?.followers || [])
	} catch {
		followingUsers.value = []
		followersUsers.value = []
	} finally {
		relationLoading.value = false
	}
}

onMounted(async () => {
	authStore.hydrate()
	await refreshUser()
	await refreshRelations()
})
</script>

<style scoped>
.profile-item {
	padding: 18px;
	border-radius: 16px;
	background: #f8fafc;
	border: 1px solid rgba(148, 163, 184, 0.2);
}

.relations-panel {
	padding-top: 14px;
}

.relations-head {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 6px;
}

.label {
	display: block;
	font-size: 12px;
	font-weight: 700;
	letter-spacing: 0.08em;
	text-transform: uppercase;
	color: #64748b;
}

.relation-tabs {
	margin-bottom: 8px;
}

.relation-tabs :deep(.el-tabs__nav-wrap::after) {
	height: 1px;
	background-color: rgba(148, 163, 184, 0.22);
}

.relation-list {
	display: grid;
	grid-template-columns: repeat(2, minmax(0, 1fr));
	gap: 10px;
}

.relation-user-card {
	display: flex;
	align-items: center;
	gap: 10px;
	padding: 10px;
	background: #ffffff;
	border-radius: 12px;
	border: 1px solid rgba(148, 163, 184, 0.16);
	cursor: pointer;
	transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.relation-actions {
	margin-left: auto;
	display: flex;
	gap: 8px;
	flex-wrap: wrap;
}

.relation-user-card:hover,
.relation-user-card:focus-visible {
	transform: translateY(-1px);
	box-shadow: 0 12px 24px rgba(15, 23, 42, 0.08);
	border-color: rgba(59, 130, 246, 0.32);
	outline: none;
}

.relation-user-meta {
	min-width: 0;
}

.relation-user-name {
	font-size: 14px;
	font-weight: 700;
	color: #0f172a;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.relation-user-email {
	margin-top: 4px;
	font-size: 12px;
	color: #64748b;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

@media (max-width: 760px) {
	.relation-list {
		grid-template-columns: 1fr;
	}

	.relation-user-card {
		flex-direction: column;
		align-items: flex-start;
	}

	.relation-actions {
		margin-left: 0;
		width: 100%;
	}
}
</style>