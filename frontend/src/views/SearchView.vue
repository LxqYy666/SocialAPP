<template>
	<div class="search-page">
		<el-card class="search-card" shadow="never">
			<div class="search-header">
				<div>
					<p class="eyebrow">Search</p>
					<h2>搜索结果</h2>
					<p class="description">
						<span v-if="keyword">“{{ keyword }}” 的搜索结果</span>
						<span v-else>请输入关键词开始搜索帖子和用户。</span>
					</p>
				</div>
				<el-button plain @click="goHome">返回首页</el-button>
			</div>

			<el-divider />

			<div v-if="!keyword" class="search-empty-state">
				<el-empty description="请输入搜索内容" :image-size="120" />
			</div>

			<template v-else>
				<div v-if="loading" class="search-skeleton">
					<el-skeleton :rows="4" animated />
				</div>

				<template v-else>
					<el-alert
						class="result-summary"
						type="info"
						:closable="false"
						:description="`找到 ${posts.length} 篇帖子和 ${users.length} 位用户`"
					/>

					<el-tabs v-model="activeTab" class="result-tabs">
						<el-tab-pane label="帖子" name="posts">
							<div v-if="posts.length" class="post-list">
								<article
									v-for="post in posts"
									:key="post.id"
									class="post-card"
									tabindex="0"
									role="button"
									@click="openPost(post.id)"
									@keydown.enter.prevent="openPost(post.id)"
									@keydown.space.prevent="openPost(post.id)"
								>
									<div class="post-top">
										<div class="post-author">
											<el-avatar :size="40">{{ authorText(post.name) }}</el-avatar>
											<div>
												<div class="post-name">{{ post.name || '未命名作者' }}</div>
												<div class="post-meta">{{ formatPostDate(post.createdAt) }}</div>
											</div>
										</div>
										<div class="post-stats">
											<el-tag size="small" type="success" effect="plain">{{ post.likes.length }} 赞</el-tag>
											<el-tag size="small" type="info">{{ post.comments.length }} 评论</el-tag>
										</div>
									</div>

									<h3 class="post-title">{{ post.title }}</h3>
									<p class="post-message">{{ post.message }}</p>

									<el-image
										v-if="post.selectedFile"
										class="post-image"
										:src="post.selectedFile"
										fit="cover"
										:preview-src-list="[post.selectedFile]"
									/>

									<div class="post-actions">
										<el-button plain @click.stop="openUserProfile(post.creator)">作者主页</el-button>
										<el-button type="primary" @click.stop="openPost(post.id)">查看详情</el-button>
									</div>
								</article>
							</div>
							<el-empty v-else description="未找到匹配帖子" :image-size="100" />
						</el-tab-pane>

						<el-tab-pane label="用户" name="users">
							<div v-if="users.length" class="user-list">
								<article v-for="user in users" :key="user.id" class="user-card">
									<div class="user-top">
										<div class="user-author">
											<el-avatar :size="48" :src="user.imageUrl || ''">{{ authorText(user.name) }}</el-avatar>
											<div>
												<div class="user-name">{{ user.name || '未命名用户' }}</div>
												<div class="user-email">{{ user.email }}</div>
											</div>
										</div>
										<el-tag size="small" type="info">{{ user.followers.length }} 粉丝</el-tag>
									</div>

									<p class="user-bio">{{ user.bio || '这个用户还没有填写简介。' }}</p>

									<div class="user-actions">
										<el-button plain @click="openUserProfile(user.id)">查看主页</el-button>
									</div>
								</article>
							</div>
							<el-empty v-else description="未找到匹配用户" :image-size="100" />
						</el-tab-pane>
					</el-tabs>
				</template>
			</template>
		</el-card>
	</div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { postApi } from '../api'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const activeTab = ref('posts')
const posts = ref([])
const users = ref([])

const keyword = computed(() => String(route.query.q || '').trim())

function normalizePosts(rawPosts = []) {
	return rawPosts.map((post) => ({
		id: post?.id || post?._id || '',
		creator: post?.creator || '',
		title: post?.title || '未命名帖子',
		message: post?.message || '',
		name: post?.name || '',
		selectedFile: post?.selectedFile || '',
		likes: Array.isArray(post?.likes) ? post.likes : [],
		comments: Array.isArray(post?.comments) ? post.comments : [],
		createdAt: post?.createdAt || '',
	}))
}

function normalizeUsers(rawUsers = []) {
	return rawUsers.map((user) => ({
		id: user?.id || user?._id || '',
		name: user?.name || '',
		email: user?.email || '',
		imageUrl: user?.imageUrl || '',
		bio: user?.bio || '',
		followers: Array.isArray(user?.followers) ? user.followers : [],
		following: Array.isArray(user?.following) ? user.following : [],
	}))
}

function authorText(name) {
	const value = String(name || '').trim()
	return value ? value.slice(0, 1).toUpperCase() : 'U'
}

function formatPostDate(createdAt) {
	if (!createdAt) {
		return '刚刚'
	}

	const date = new Date(createdAt)
	if (Number.isNaN(date.getTime())) {
		return '刚刚'
	}

	return new Intl.DateTimeFormat('zh-CN', {
		year: 'numeric',
		month: '2-digit',
		day: '2-digit',
		hour: '2-digit',
		minute: '2-digit',
	}).format(date)
}

async function loadSearchResults() {
	const query = keyword.value
	if (!query) {
		posts.value = []
		users.value = []
		return
	}

	loading.value = true
	try {
		const response = await postApi.search(query)
		posts.value = normalizePosts(response?.posts || [])
		users.value = normalizeUsers(response?.users || [])
		if (posts.value.length) {
			activeTab.value = 'posts'
		} else if (users.value.length) {
			activeTab.value = 'users'
		}
	} catch {
		posts.value = []
		users.value = []
	} finally {
		loading.value = false
	}
}

function openPost(postId) {
	if (!postId) {
		return
	}

	router.push({ name: 'post-detail', params: { id: postId } })
}

function openUserProfile(userId) {
	if (!userId) {
		return
	}

	router.push({ name: 'user-profile', params: { id: userId } })
}

function goHome() {
	router.push({ name: 'home' })
}

onMounted(() => {
	authStore.hydrate()
})

watch(keyword, loadSearchResults, { immediate: true })
</script>

<style scoped>
.search-page {
	min-height: calc(100vh - 72px);
	padding: 28px 20px;
	background:
		radial-gradient(circle at top right, rgba(37, 99, 235, 0.12), transparent 36%),
		radial-gradient(circle at bottom left, rgba(16, 185, 129, 0.12), transparent 32%),
		linear-gradient(180deg, #f7fbff 0%, #eef6ff 100%);
}

.search-card {
	max-width: 1080px;
	margin: 0 auto;
	border-radius: 24px;
	border: 1px solid rgba(148, 163, 184, 0.14);
	box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08);
}

.search-header {
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

.search-header h2 {
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

.search-skeleton,
.search-empty-state {
	padding: 16px 0 0;
}

.result-summary {
	margin-bottom: 16px;
}

.result-tabs :deep(.el-tabs__header) {
	margin-bottom: 18px;
}

.post-list,
.user-list {
	display: flex;
	flex-direction: column;
	gap: 14px;
}

.post-card,
.user-card {
	padding: 18px;
	border-radius: 18px;
	background: #ffffff;
	border: 1px solid rgba(148, 163, 184, 0.16);
	box-shadow: 0 10px 24px rgba(15, 23, 42, 0.04);
	transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.post-card:hover,
.user-card:hover {
	transform: translateY(-1px);
	box-shadow: 0 14px 30px rgba(15, 23, 42, 0.08);
}

.post-top,
.user-top {
	display: flex;
	align-items: flex-start;
	justify-content: space-between;
	gap: 12px;
}

.post-author,
.user-author {
	display: flex;
	align-items: center;
	gap: 12px;
}

.post-name,
.user-name {
	font-size: 15px;
	font-weight: 700;
	color: #0f172a;
}

.post-meta,
.user-email {
	margin-top: 4px;
	font-size: 12px;
	color: #64748b;
}

.post-stats {
	display: flex;
	gap: 8px;
	flex-wrap: wrap;
}

.post-title {
	margin: 14px 0 8px;
	font-size: 18px;
	line-height: 1.35;
	color: #0f172a;
}

.post-message,
.user-bio {
	margin: 0;
	color: #334155;
	line-height: 1.8;
	white-space: pre-wrap;
}

.post-image {
	display: block;
	width: 100%;
	max-height: 420px;
	margin-top: 14px;
	border-radius: 14px;
	overflow: hidden;
}

.post-actions,
.user-actions {
	margin-top: 16px;
	display: flex;
	justify-content: flex-end;
	gap: 10px;
	flex-wrap: wrap;
}

@media (max-width: 640px) {
	.search-page {
		padding: 16px 12px;
	}

	.search-header {
		flex-direction: column;
		align-items: flex-start;
	}

	.search-header h2 {
		font-size: 26px;
	}

	.post-top,
	.user-top {
		flex-direction: column;
	}

	.post-actions,
	.user-actions {
		justify-content: flex-start;
	}
}
</style>