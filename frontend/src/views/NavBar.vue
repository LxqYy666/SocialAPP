<template>
	<header class="nav-bar">
		<el-row class="nav-row" align="middle" justify="space-between" :gutter="16">
			<el-col :xs="24" :sm="6" :md="5" class="nav-left">
				<el-link class="brand-link" href="/" :underline="false" aria-label="返回主页">
					<el-avatar :size="42" shape="square" src="/vite.svg" />
				</el-link>
			</el-col>

			<el-col :xs="24" :sm="12" :md="12" class="nav-center">
				<el-input
					v-model="searchText"
					class="search-input"
					placeholder="搜索内容"
					clearable
					size="large"
					@keyup.enter="submitSearch"
					@clear="submitSearch"
				>
					<template #prefix>
						<svg viewBox="0 0 24 24" aria-hidden="true" class="search-icon">
							<path
								d="M10.5 4a6.5 6.5 0 1 1 0 13 6.5 6.5 0 0 1 0-13Zm0 2a4.5 4.5 0 1 0 0 9 4.5 4.5 0 0 0 0-9Zm5.97 8.56 3.97 3.97-1.41 1.41-3.97-3.97 1.41-1.41Z"
							/>
						</svg>
					</template>
					<template #suffix>
						<el-button class="search-action" circle text aria-label="执行搜索" @click="submitSearch">
							<svg viewBox="0 0 24 24" aria-hidden="true" class="search-icon">
								<path
									d="M10.5 4a6.5 6.5 0 1 1 0 13 6.5 6.5 0 0 1 0-13Zm0 2a4.5 4.5 0 1 0 0 9 4.5 4.5 0 0 0 0-9Zm5.97 8.56 3.97 3.97-1.41 1.41-3.97-3.97 1.41-1.41Z"
								/>
							</svg>
						</el-button>
					</template>
				</el-input>
			</el-col>

			<el-col :xs="24" :sm="6" :md="7" class="nav-right">
				<el-space :size="10" alignment="center">
					<el-button type="primary" class="publish-button" @click="openPublishDialog">
						发布帖子
					</el-button>

					<el-badge :value="unreadMessageCount" :hidden="unreadMessageCount === 0" class="notification-badge">
						<el-button class="icon-button" circle aria-label="消息" @click="openMessages">
							<svg viewBox="0 0 24 24" aria-hidden="true" class="action-icon">
								<path
									d="M4 4h16v12H7.17L4 19.17V4Zm2 2v8.34L6.34 14H18V6H6Zm2 2h8v2H8V8Zm0 4h6v2H8v-2Z"
								/>
							</svg>
						</el-button>
					</el-badge>

					<el-popover
						placement="bottom-end"
						:width="380"
						trigger="click"
						:show-after="0"
						:hide-after="80"
						@show="loadNotifications"
					>
						<template #reference>
							<el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notification-badge">
								<el-button class="icon-button" circle aria-label="通知">
									<svg viewBox="0 0 24 24" aria-hidden="true" class="action-icon">
										<path
											d="M12 3a5 5 0 0 0-5 5v2.28c0 .79-.25 1.56-.72 2.2L4 15v1h16v-1l-2.28-2.52c-.47-.64-.72-1.41-.72-2.2V8a5 5 0 0 0-5-5Zm0 18a2.5 2.5 0 0 0 2.45-2h-4.9A2.5 2.5 0 0 0 12 21Z"
										/>
									</svg>
								</el-button>
							</el-badge>
						</template>

						<div class="notification-popover">
							<div class="notification-head">
								<div>
									<div class="notification-title">通知</div>
									<div class="notification-subtitle">{{ unreadCount }} 条未读</div>
								</div>
								<el-button size="small" plain :disabled="!unreadCount" @click="markAllNotificationsRead">
									全部已读
								</el-button>
							</div>

							<div v-if="notificationLoading" class="notification-state">
								<el-skeleton :rows="3" animated />
							</div>

							<div v-else-if="notifications.length" class="notification-list">
								<div
									v-for="notification in notifications"
									:key="notification.id"
									class="notification-item"
									:class="{ unread: !notification.isReaded }"
									tabindex="0"
									role="button"
									@click="openNotification(notification)"
									@keydown.enter.prevent="openNotification(notification)"
									@keydown.space.prevent="openNotification(notification)"
								>
									<el-avatar :size="40" :src="notification.notificationUser.avatar || ''">
										{{ notificationUserText(notification) }}
									</el-avatar>
									<div class="notification-content">
										<div class="notification-text">{{ notification.details }}</div>
										<div class="notification-meta">
											<span>{{ notification.notificationUser.name || '系统通知' }}</span>
											<span>{{ formatNotificationDate(notification.createdAt) }}</span>
										</div>
									</div>
									<el-tag size="small" :type="notification.isReaded ? 'info' : 'danger'" effect="plain">
										{{ notification.isReaded ? '已读' : '未读' }}
									</el-tag>
								</div>
							</div>

							<el-empty v-else description="暂无通知" :image-size="84" />
						</div>
					</el-popover>

					<el-popover
						placement="bottom-end"
						:width="280"
						trigger="hover"
						:show-after="120"
						:hide-after="80"
					>
						<template #reference>
							<el-avatar class="avatar-button" :size="42" :src="profileImage" aria-label="个人头像">
								{{ avatarText }}
							</el-avatar>
						</template>

						<div class="profile-popover">
							<div class="profile-actions">
								<el-button size="small" plain @click="goProfile">个人信息</el-button>
								<el-button size="small" type="danger" @click="handleLogout">登出</el-button>
							</div>
						</div>
					</el-popover>
				</el-space>
			</el-col>
		</el-row>

		<el-dialog v-model="publishDialogVisible" title="发布帖子" width="min(720px, 92vw)">
			<el-form ref="postFormRef" :model="postForm" :rules="postRules" label-position="top" class="publish-form">
				<el-form-item label="帖子标题" prop="title">
					<el-input v-model="postForm.title" maxlength="80" show-word-limit placeholder="输入帖子标题" />
				</el-form-item>

				<el-form-item label="帖子内容" prop="message">
					<el-input
						v-model="postForm.message"
						type="textarea"
						:rows="6"
						maxlength="1000"
						show-word-limit
						placeholder="分享一些动态、想法或故事"
					/>
				</el-form-item>

				<el-form-item label="图片地址 / Base64 编码（可选）" prop="selectedFile">
					<el-input
						v-model="postForm.selectedFile"
						type="textarea"
						:rows="4"
						maxlength="10000"
						show-word-limit
						placeholder="粘贴图片 URL 或 base64 字符串"
					/>
					<div class="helper-text">如果不需要图片，可以留空。</div>
				</el-form-item>
			</el-form>

			<template #footer>
				<div class="dialog-actions">
					<el-button @click="closePublishDialog">取消</el-button>
					<el-button :disabled="publishing" @click="resetPostForm">清空</el-button>
					<el-button type="primary" :loading="publishing" @click="submitPost">发布帖子</el-button>
				</div>
			</template>
		</el-dialog>
	</header>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'

import { chatApi, postApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { useNotificationRealtimeStore } from '../stores/notificationRealtime'

const searchText = ref('')
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const notificationStore = useNotificationRealtimeStore()
const publishDialogVisible = ref(false)
const publishing = ref(false)
const postFormRef = ref()
const unreadMessageCount = ref(0)

const profileImage = computed(() => authStore.user?.imageUrl || '')
const userId = computed(() => authStore.userId || '')
const notificationLoading = computed(() => notificationStore.loading)
const notifications = computed(() => notificationStore.notifications)
const unreadCount = computed(() => notificationStore.unreadCount)

const avatarText = computed(() => {
	const name = (authStore.user?.name || '').trim()
	return name ? name.slice(0, 1).toUpperCase() : 'U'
})

const postForm = reactive({
	title: '',
	message: '',
	selectedFile: '',
})

const postRules = {
	title: [{ required: true, message: '请输入帖子标题', trigger: 'blur' }],
	message: [{ required: true, message: '请输入帖子内容', trigger: 'blur' }],
}

function submitSearch() {
	const keyword = searchText.value.trim()
	if (!keyword) {
		ElMessage.warning('请输入搜索内容')
		return
	}

	router.push({ name: 'search', query: { q: keyword } })
}

async function loadUnreadMessageCount() {
	if (!userId.value) {
		unreadMessageCount.value = 0
		return
	}

	try {
		const response = await chatApi.getUnreadMessages()
		unreadMessageCount.value = Number(response?.totalUnReadedMsg || 0)
	} catch {
		unreadMessageCount.value = 0
	}
}

function notificationUserText(notification) {
	const name = String(notification?.notificationUser?.name || '').trim()
	return name ? name.slice(0, 1).toUpperCase() : 'N'
}

function formatNotificationDate(createdAt) {
	if (!createdAt) {
		return '刚刚'
	}

	const date = new Date(createdAt)
	if (Number.isNaN(date.getTime())) {
		return '刚刚'
	}

	return new Intl.DateTimeFormat('zh-CN', {
		month: '2-digit',
		day: '2-digit',
		hour: '2-digit',
		minute: '2-digit',
	}).format(date)
}

async function loadNotifications() {
	try {
		await notificationStore.fetchNotifications()
	} catch {
		notificationStore.clear()
	}
}

async function markAllNotificationsRead() {
	if (!unreadCount.value) {
		return
	}

	try {
		await notificationStore.markAllRead()
		ElMessage.success('已标记全部通知为已读')
	} catch (error) {
		ElMessage.error(error?.response?.data?.error || '标记已读失败')
	}
}

async function openNotification(notification) {
	if (!notification?.id) {
		return
	}

	if (!notification.isReaded) {
		try {
			await notificationStore.markAllRead()
		} catch {
			// keep navigation usable even if the read update fails
		}
	}

	if (notification.type === 'follow') {
		router.push({ name: 'user-profile', params: { id: notification.targetUserId } })
		return
	}

	router.push({ name: 'post-detail', params: { id: notification.targetUserId } })
}

function openPublishDialog() {
	publishDialogVisible.value = true
}

function closePublishDialog() {
	publishDialogVisible.value = false
}

function resetPostForm() {
	postForm.title = ''
	postForm.message = ''
	postForm.selectedFile = ''
	postFormRef.value?.clearValidate()
}

async function submitPost() {
	try {
		await postFormRef.value?.validate()
	} catch {
		return
	}

	publishing.value = true
	try {
		await postApi.create({
			title: postForm.title.trim(),
			message: postForm.message.trim(),
			selectedFile: postForm.selectedFile.trim(),
		})
		ElMessage.success('帖子发布成功')
		resetPostForm()
		closePublishDialog()
	} catch (error) {
		ElMessage.error(error?.response?.data?.error || '发布失败，请稍后重试')
	} finally {
		publishing.value = false
	}
}

function goProfile() {
	router.push({ name: 'profile' })
}

function handleLogout() {
	notificationStore.disconnect()
	notificationStore.clear()
	authStore.logout()
	router.replace({ name: 'auth', params: { mode: 'login' } })
}

function openMessages() {
	router.push({ name: 'messages' })
}

watch(
	() => route.query.q,
	(query) => {
		searchText.value = typeof query === 'string' ? query : ''
	},
	{ immediate: true },
)

onMounted(() => {
	authStore.hydrate()
	loadUnreadMessageCount()
	if (userId.value) {
		notificationStore.connect(userId.value)
	}
	window.addEventListener('chat-unread-changed', loadUnreadMessageCount)
})

onBeforeUnmount(() => {
	notificationStore.disconnect()
	window.removeEventListener('chat-unread-changed', loadUnreadMessageCount)
})

watch(
	userId,
	(nextUserId) => {
		if (!nextUserId) {
			notificationStore.disconnect()
			notificationStore.clear()
			unreadMessageCount.value = 0
			return
		}

		notificationStore.connect(nextUserId)
		loadNotifications()
		loadUnreadMessageCount()
	},
	{ immediate: true },
)
</script>

<style scoped>
.nav-bar {
	padding: 10px 24px;
	background: linear-gradient(135deg, #ffffff 0%, #f7f9fc 100%);
	border-bottom: 1px solid rgba(15, 23, 42, 0.08);
	box-shadow: 0 10px 30px rgba(15, 23, 42, 0.04);
}

.nav-row {
	width: 100%;
	min-height: 72px;
}

.nav-left,
.nav-center,
.nav-right {
	display: flex;
	align-items: center;
	min-width: 0;
}

.nav-left {
	justify-content: flex-start;
}

.nav-center {
	justify-content: center;
}

.nav-right {
	justify-content: flex-end;
}

.brand-link {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	text-decoration: none;
}

.search-input {
	width: 100%;
	max-width: 560px;
}

.search-input :deep(.el-input__wrapper) {
	border-radius: 999px;
	background: #eef2f7;
	box-shadow: none;
	padding-left: 14px;
	padding-right: 14px;
	transition: border-color 0.2s ease, background-color 0.2s ease, box-shadow 0.2s ease;
}

.search-input :deep(.el-input__wrapper.is-focus) {
	background: #ffffff;
	box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.08);
}

.search-input :deep(.el-input__prefix) {
	margin-right: 8px;
}

.search-input :deep(.el-input__inner) {
	color: #0f172a;
	font-size: 14px;
}

.search-input :deep(.el-input__inner::placeholder) {
	color: #94a3b8;
}

.search-icon,
.action-icon {
	width: 18px;
	height: 18px;
	fill: currentColor;
	display: block;
}

.icon-button {
	width: 42px;
	height: 42px;
	padding: 0;
	background: #eef2f7;
	border: 0;
	border-radius: 12px;
	color: #334155;
	transition: transform 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
}

.notification-badge {
	display: inline-flex;
}

.publish-button {
	height: 42px;
	padding: 0 14px;
	border-radius: 12px;
	box-shadow: none;
}

.avatar-button {
	width: 42px;
	height: 42px;
	padding: 0;
	background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
	cursor: pointer;
	border: 0;
	color: #334155;
	transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.icon-button:hover,
.avatar-button:hover {
	transform: translateY(-1px);
	box-shadow: 0 10px 20px rgba(15, 23, 42, 0.08);
}

.avatar-button :deep(img) {
	width: 100%;
	height: 100%;
	object-fit: cover;
	border-radius: 50%;
}

.profile-popover {
	padding: 4px;
}

.profile-actions {
	display: flex;
	justify-content: flex-end;
	gap: 8px;
}

.publish-form {
	padding-top: 4px;
}

.helper-text {
	margin-top: 8px;
	font-size: 12px;
	color: #64748b;
}

.dialog-actions {
	display: flex;
	justify-content: flex-end;
	gap: 10px;
}

.notification-popover {
	display: flex;
	flex-direction: column;
	gap: 12px;
	max-height: 480px;
}

.notification-head {
	display: flex;
	align-items: flex-start;
	justify-content: space-between;
	gap: 12px;
	padding-bottom: 10px;
	border-bottom: 1px solid rgba(148, 163, 184, 0.18);
}

.notification-title {
	font-size: 16px;
	font-weight: 800;
	color: #0f172a;
}

.notification-subtitle {
	margin-top: 4px;
	font-size: 12px;
	color: #64748b;
}

.notification-state {
	padding: 12px 0;
}

.notification-list {
	display: flex;
	flex-direction: column;
	gap: 10px;
	max-height: 360px;
	overflow: auto;
}

.notification-item {
	display: grid;
	grid-template-columns: auto 1fr auto;
	gap: 12px;
	align-items: center;
	padding: 12px;
	border-radius: 14px;
	background: #f8fafc;
	border: 1px solid rgba(148, 163, 184, 0.12);
	cursor: pointer;
	transition: background-color 0.2s ease, transform 0.2s ease;
}

.notification-item:hover {
	background: #eef4ff;
	transform: translateY(-1px);
}

.notification-item.unread {
	background: rgba(59, 130, 246, 0.08);
	border-color: rgba(59, 130, 246, 0.18);
}

.notification-content {
	min-width: 0;
}

.notification-text {
	font-size: 13px;
	font-weight: 600;
	color: #0f172a;
	line-height: 1.5;
	word-break: break-word;
}

.notification-meta {
	margin-top: 5px;
	display: flex;
	flex-wrap: wrap;
	gap: 8px;
	font-size: 12px;
	color: #64748b;
}

@media (max-width: 900px) {
	.nav-bar {
		padding: 10px 16px;
	}

	.nav-center {
		justify-content: flex-start;
	}
}

@media (max-width: 640px) {
	.nav-left,
	.nav-center,
	.nav-right {
		justify-content: center;
	}

	.nav-center {
		margin: 6px 0;
	}

	.nav-right {
		justify-content: center;
	}
}
</style>
