<template>
	<div class="chat-page">
		<el-card class="chat-card" shadow="never">
			<div class="chat-header">
				<div>
					<p class="eyebrow">Messages</p>
					<h2>聊天</h2>
					<p class="description">和你的关注人、粉丝继续交流。</p>
				</div>
				<el-button plain @click="goHome">返回首页</el-button>
			</div>

			<el-divider />

			<div class="chat-layout">
				<aside class="contact-panel">
					<div class="panel-head">
						<span>联系人</span>
						<el-tag size="small" type="info">{{ contacts.length }}</el-tag>
					</div>

					<div v-if="contactsLoading" class="panel-state">
						<el-skeleton :rows="4" animated />
					</div>

					<div v-else-if="contacts.length" class="contact-list">
						<div
							v-for="contact in contacts"
							:key="contact.id"
							class="contact-item"
							:class="{ active: activeContactId === contact.id }"
							tabindex="0"
							role="button"
							@click="selectContact(contact.id)"
							@keydown.enter.prevent="selectContact(contact.id)"
							@keydown.space.prevent="selectContact(contact.id)"
						>
							<el-avatar :size="44" :src="contact.imageUrl || ''">{{ authorText(contact.name) }}</el-avatar>
							<div class="contact-meta">
								<div class="contact-name">{{ contact.name || '未命名用户' }}</div>
								<div class="contact-tags">
									<el-tag v-if="contact.relation === 'mutual'" size="small" type="success">互相关注</el-tag>
									<el-tag v-else-if="contact.relation === 'following'" size="small" type="info">已关注</el-tag>
									<el-tag v-else size="small" type="warning">粉丝</el-tag>
									<el-tag v-if="contact.unreadCount" size="small" type="danger">{{ contact.unreadCount }}</el-tag>
								</div>
							</div>
						</div>
					</div>
					<el-empty v-else description="暂无联系人" :image-size="90" />
				</aside>

				<section class="conversation-panel">
					<template v-if="activeContact">
						<div class="conversation-head">
							<div class="conversation-user">
								<el-avatar :size="46" :src="activeContact.imageUrl || ''">{{ authorText(activeContact.name) }}</el-avatar>
								<div>
									<div class="conversation-name">{{ activeContact.name }}</div>
									<div class="conversation-subtitle">{{ activeContact.email || '暂无邮箱' }}</div>
								</div>
							</div>
							<div class="conversation-actions">
								<el-button plain @click="openProfile(activeContact.id)">查看主页</el-button>
								<el-button type="primary" plain @click="markCurrentChatRead">标记已读</el-button>
							</div>
						</div>

						<div ref="messageListRef" class="message-list">
							<div v-if="messageLoading && !messages.length" class="panel-state">
								<el-skeleton :rows="5" animated />
							</div>

							<div v-else-if="messages.length">
								<div class="load-more-wrap">
									<el-button v-if="hasMoreHistory" size="small" plain :loading="historyLoading" @click="loadMoreHistory">
										加载更早消息
									</el-button>
								</div>

								<div v-for="(message, index) in messages" :key="`${message.id || index}-${index}`" class="message-row" :class="message.sender === currentUserId ? 'mine' : 'theirs'">
									<div class="message-bubble">
										{{ message.content }}
									</div>
								</div>
							</div>

							<el-empty v-else description="开始发送第一条消息吧" :image-size="96" />
						</div>

						<div class="composer">
							<el-input
								v-model="messageText"
								type="textarea"
								:rows="3"
								maxlength="1000"
								show-word-limit
								placeholder="输入消息，回车发送"
								@keydown.enter.exact.prevent="sendMessage"
							/>
							<div class="composer-actions">
								<el-button :disabled="sending" @click="messageText = ''">清空</el-button>
								<el-button type="primary" :loading="sending" @click="sendMessage">发送</el-button>
							</div>
						</div>
					</template>

					<el-empty v-else description="请选择一个联系人开始聊天" :image-size="120" />
				</section>
			</div>
		</el-card>
	</div>
</template>

<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'

import { chatApi, userApi } from '../api'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const contactsLoading = ref(false)
const messageLoading = ref(false)
const historyLoading = ref(false)
const sending = ref(false)
const contacts = ref([])
const messages = ref([])
const messageText = ref('')
const activeContactId = ref('')
const historyPage = ref(0)
const hasMoreHistory = ref(true)
const messageListRef = ref(null)

const currentUserId = computed(() => authStore.userId || '')
const activeContact = computed(() => contacts.value.find((contact) => contact.id === activeContactId.value) || null)

function authorText(name) {
	const value = String(name || '').trim()
	return value ? value.slice(0, 1).toUpperCase() : 'U'
}

function normalizeContactUser(rawUser, relation) {
	return {
		id: rawUser?.id || rawUser?._id || '',
		name: rawUser?.name || '未命名用户',
		email: rawUser?.email || '',
		imageUrl: rawUser?.imageUrl || '',
		relation,
		unreadCount: 0,
	}
}

function uniqueById(list = []) {
	const map = new Map()
	list.forEach((item) => {
		if (item?.id && !map.has(item.id)) {
			map.set(item.id, item)
		}
	})
	return Array.from(map.values())
}

function normalizeMessages(rawMessages = []) {
	return rawMessages.map((message) => ({
		id: message?.id || message?._id || '',
		sender: message?.sender || '',
		receiver: message?.receiver || '',
		content: message?.content || '',
	}))
}

async function fetchContacts() {
	if (!currentUserId.value) {
		contacts.value = []
		return
	}

	contactsLoading.value = true
	try {
		const response = await userApi.getById(currentUserId.value)
		const user = response?.user || authStore.user || {}
		const followingIds = Array.isArray(user.following) ? user.following : []
		const followersIds = Array.isArray(user.followers) ? user.followers : []
		const followingUsers = await Promise.allSettled(followingIds.map((id) => userApi.getById(id)))
		const followerUsers = await Promise.allSettled(followersIds.map((id) => userApi.getById(id)))
		const followingList = followingUsers
			.filter((item) => item.status === 'fulfilled' && item.value?.user)
			.map((item) => normalizeContactUser(item.value.user, 'following'))
		const followerList = followerUsers
			.filter((item) => item.status === 'fulfilled' && item.value?.user)
			.map((item) => normalizeContactUser(item.value.user, 'follower'))
		const combined = uniqueById([...followingList, ...followerList])
		const followingSet = new Set(followingIds.map(String))
		const followerSet = new Set(followersIds.map(String))
		contacts.value = combined.map((contact) => ({
			...contact,
			relation: followingSet.has(contact.id) && followerSet.has(contact.id)
				? 'mutual'
				: followingSet.has(contact.id)
					? 'following'
					: 'follower',
		}))
		if (!activeContactId.value && contacts.value.length) {
			activeContactId.value = contacts.value[0].id
		}
	} catch {
		contacts.value = []
	} finally {
		contactsLoading.value = false
	}
}

async function refreshUnreadCounts() {
	if (!currentUserId.value) {
		return
	}

	try {
		const response = await chatApi.getUnreadMessages()
		const unreadMap = new Map(
			(response?.unReadedMsgs || []).map((item) => [item.otherUserId, Number(item.numOfUnReadedMsg || 0)]),
		)
		contacts.value = contacts.value.map((contact) => ({
			...contact,
			unreadCount: unreadMap.get(contact.id) || 0,
		}))
		window.dispatchEvent(new Event('chat-unread-changed'))
	} catch {
		// keep current unread badges if the refresh fails
	}
}

async function loadMessages(page = 0, replace = false) {
	if (!activeContactId.value || !currentUserId.value) {
		messages.value = []
		return
	}

	if (replace) {
		messageLoading.value = true
	} else {
		historyLoading.value = true
	}

	try {
		const response = await chatApi.getMessages({
			from: page,
			firstuid: currentUserId.value,
			seconduid: activeContactId.value,
		})
		const nextMessages = normalizeMessages(response?.messages || [])
		messages.value = replace ? nextMessages : [...nextMessages, ...messages.value]
		historyPage.value = page
		hasMoreHistory.value = nextMessages.length > 0
		await nextTick()
		scrollToBottom()
	} catch {
		if (replace) {
			messages.value = []
		}
	} finally {
		messageLoading.value = false
		historyLoading.value = false
	}
}

function scrollToBottom() {
	const container = messageListRef.value
	if (container) {
		container.scrollTop = container.scrollHeight
	}
}

async function selectContact(contactId) {
	if (!contactId || contactId === activeContactId.value) {
		return
	}

	activeContactId.value = contactId
	historyPage.value = 0
	hasMoreHistory.value = true
	await loadMessages(0, true)
	await markCurrentChatRead()
	await refreshUnreadCounts()
}

async function loadMoreHistory() {
	if (!hasMoreHistory.value || historyLoading.value) {
		return
	}

	await loadMessages(historyPage.value + 1, false)
}

async function markCurrentChatRead() {
	if (!activeContactId.value) {
		return
	}

	try {
		await chatApi.markUnreadMessages(activeContactId.value)
		await refreshUnreadCounts()
	} catch {
		// ignore read-state failure so chat remains usable
	}
}

async function sendMessage() {
	const content = messageText.value.trim()
	if (!content || !activeContactId.value) {
		return
	}

	sending.value = true
	try {
		await chatApi.sendMessage({
			receiver: activeContactId.value,
			content,
		})
		messageText.value = ''
		await loadMessages(0, true)
		await refreshUnreadCounts()
		window.dispatchEvent(new Event('chat-unread-changed'))
	} catch (error) {
		ElMessage.error(error?.response?.data?.error || '发送失败，请稍后重试')
	} finally {
		sending.value = false
	}
}

function openProfile(userId) {
	if (!userId) {
		return
	}

	router.push({ name: 'user-profile', params: { id: userId } })
}

function goHome() {
	router.push({ name: 'home' })
}

onMounted(async () => {
	authStore.hydrate()
	await fetchContacts()
	await refreshUnreadCounts()
	const targetId = String(route.query.target || '')
	if (targetId) {
		activeContactId.value = targetId
		await loadMessages(0, true)
		await markCurrentChatRead()
	}
})

watch(
	() => route.query.target,
	async (target) => {
		const nextTarget = String(target || '')
		if (!nextTarget || nextTarget === activeContactId.value) {
			return
		}

		activeContactId.value = nextTarget
		historyPage.value = 0
		hasMoreHistory.value = true
		await loadMessages(0, true)
		await markCurrentChatRead()
	},
)

watch(activeContactId, async () => {
	await refreshUnreadCounts()
})
</script>

<style scoped>
.chat-page {
	min-height: calc(100vh - 72px);
	padding: 28px 20px;
	background:
		radial-gradient(circle at top right, rgba(37, 99, 235, 0.12), transparent 36%),
		radial-gradient(circle at bottom left, rgba(16, 185, 129, 0.12), transparent 32%),
		linear-gradient(180deg, #f7fbff 0%, #eef6ff 100%);
}

.chat-card {
	max-width: 1240px;
	margin: 0 auto;
	border-radius: 24px;
	border: 1px solid rgba(148, 163, 184, 0.14);
	box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08);
}

.chat-header {
	display: flex;
	align-items: flex-end;
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

.chat-header h2 {
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

.chat-layout {
	display: grid;
	grid-template-columns: 320px minmax(0, 1fr);
	gap: 16px;
}

.contact-panel,
.conversation-panel {
	min-width: 0;
	padding: 16px;
	border-radius: 18px;
	background: #ffffff;
	border: 1px solid rgba(148, 163, 184, 0.16);
	box-shadow: 0 10px 24px rgba(15, 23, 42, 0.04);
}

.panel-head {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 12px;
	font-weight: 700;
	color: #0f172a;
}

.panel-state {
	padding: 12px 0;
}

.contact-list {
	display: flex;
	flex-direction: column;
	gap: 10px;
	max-height: 620px;
	overflow: auto;
}

.contact-item {
	display: flex;
	align-items: center;
	gap: 12px;
	padding: 12px;
	border-radius: 14px;
	border: 1px solid rgba(148, 163, 184, 0.16);
	cursor: pointer;
	transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.contact-item:hover,
.contact-item:focus-visible,
.contact-item.active {
	transform: translateY(-1px);
	box-shadow: 0 12px 24px rgba(15, 23, 42, 0.08);
	border-color: rgba(59, 130, 246, 0.28);
	outline: none;
}

.contact-meta {
	min-width: 0;
	flex: 1;
}

.contact-name,
.conversation-name {
	font-size: 15px;
	font-weight: 700;
	color: #0f172a;
}

.contact-tags {
	margin-top: 6px;
	display: flex;
	flex-wrap: wrap;
	gap: 6px;
}

.conversation-panel {
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.conversation-head {
	display: flex;
	align-items: flex-start;
	justify-content: space-between;
	gap: 12px;
	padding-bottom: 12px;
	border-bottom: 1px solid rgba(148, 163, 184, 0.16);
}

.conversation-user {
	display: flex;
	align-items: center;
	gap: 12px;
}

.conversation-subtitle {
	margin-top: 4px;
	font-size: 12px;
	color: #64748b;
}

.conversation-actions {
	display: flex;
	gap: 8px;
	flex-wrap: wrap;
}

.message-list {
	min-height: 520px;
	max-height: 620px;
	overflow: auto;
	display: flex;
	flex-direction: column;
	gap: 10px;
	padding-right: 4px;
}

.load-more-wrap {
	display: flex;
	justify-content: center;
	padding-bottom: 6px;
}

.message-row {
	display: flex;
}

.message-row.mine {
	justify-content: flex-end;
}

.message-row.theirs {
	justify-content: flex-start;
}

.message-bubble {
	max-width: min(72%, 520px);
	padding: 12px 14px;
	border-radius: 16px;
	line-height: 1.7;
	white-space: pre-wrap;
	word-break: break-word;
	background: #f8fafc;
	color: #0f172a;
	border: 1px solid rgba(148, 163, 184, 0.16);
}

.message-row.mine .message-bubble {
	background: linear-gradient(135deg, rgba(37, 99, 235, 0.14), rgba(16, 185, 129, 0.14));
	border-color: rgba(37, 99, 235, 0.16);
}

.composer {
	display: flex;
	flex-direction: column;
	gap: 10px;
	padding-top: 8px;
	border-top: 1px solid rgba(148, 163, 184, 0.16);
}

.composer-actions {
	display: flex;
	justify-content: flex-end;
	gap: 10px;
}

@media (max-width: 960px) {
	.chat-layout {
		grid-template-columns: 1fr;
	}
}

@media (max-width: 640px) {
	.chat-page {
		padding: 16px 12px;
	}

	.chat-header {
		flex-direction: column;
		align-items: flex-start;
	}

	.chat-header h2 {
		font-size: 26px;
	}

	.conversation-head {
		flex-direction: column;
	}

	.conversation-actions {
		width: 100%;
	}

	.conversation-actions .el-button {
		flex: 1;
	}

	.message-bubble {
		max-width: 90%;
	}
}
</style>