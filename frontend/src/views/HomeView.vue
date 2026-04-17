<template>
	<div class="home-page">
		<div class="home-layout">
			<aside class="home-sidebar">
				<div class="sidebar-brand">
					<div class="brand-mark">S</div>
					<div>
						<div class="brand-name">Social APP</div>
						<div class="brand-subtitle">内容推荐中心</div>
					</div>
				</div>

				<el-menu
					:default-active="activeSection"
					class="section-menu"
					@select="handleSectionSelect"
				>
					<el-menu-item index="posts">
						<span>推荐帖子</span>
					</el-menu-item>
					<el-menu-item index="friends">
						<span>推荐朋友</span>
					</el-menu-item>
				</el-menu>
			</aside>

			<main class="home-main">
				<el-card class="home-card" shadow="never">
					<div class="home-header">
						<div class="header-copy">
							<p class="eyebrow">Home</p>
							<h2>{{ sectionTitle }}</h2>
							<p class="description">{{ sectionDescription }}</p>
						</div>
					</div>

					<el-divider />

					<template v-if="isPostsSection">
						<div v-if="postLoading && !posts.length" class="feed-skeleton">
							<el-skeleton :rows="4" animated />
						</div>

						<div v-else-if="posts.length" class="feed-list">
							<article v-for="post in posts" :key="post.id" class="feed-card">
								<div class="feed-top">
									<div class="feed-author">
										<el-avatar :size="40">{{ authorText(post.name) }}</el-avatar>
										<div>
											<div class="feed-name">{{ post.name || '未命名作者' }}</div>
											<div class="feed-meta">{{ formatPostDate(post.createdAt) }}</div>
										</div>
									</div>
									<div class="feed-counts">
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

								<div class="feed-actions">
									<el-button plain @click="openUserProfile(post.creator)">作者主页</el-button>
									<el-button type="primary" @click="openPostDetail(post.id)">查看详情</el-button>
								</div>
							</article>
						</div>

						<el-empty v-else description="暂无推荐帖子" :image-size="110">
							<template #description>
								<div class="empty-text">暂无推荐帖子，先去关注一些用户吧。</div>
							</template>
						</el-empty>

						<div v-if="numberOfPages > 1" class="pagination-wrap">
							<el-pagination
								v-model:current-page="currentPage"
								background
								layout="prev, pager, next"
								:total="numberOfPages * pageSize"
								:page-size="pageSize"
								:disabled="postLoading"
								@current-change="handlePageChange"
							/>
						</div>
					</template>

					<template v-else>
						<div v-if="friendLoading && !suggestedUsers.length" class="feed-skeleton">
							<el-skeleton :rows="4" animated />
						</div>

						<div v-else-if="suggestedUsers.length" class="friend-list">
							<article v-for="user in suggestedUsers" :key="user.id" class="friend-card">
								<div class="friend-top">
									<div class="friend-author">
										<el-avatar :size="48" :src="user.imageUrl || ''">{{ authorText(user.name) }}</el-avatar>
										<div>
											<div class="friend-name">{{ user.name || '未命名用户' }}</div>
											<div class="friend-email">{{ user.email }}</div>
										</div>
									</div>
									<el-tag v-if="isFollowing(user.id)" size="small" type="success">已关注</el-tag>
								</div>

								<p class="friend-bio">{{ user.bio || '这个用户还没有填写简介。' }}</p>

								<div class="friend-stats">
									<el-tag size="small" type="info">{{ user.followers.length }} 粉丝</el-tag>
									<el-tag size="small" type="info">{{ user.following.length }} 关注</el-tag>
								</div>

								<div class="friend-actions">
									<el-button plain @click="openUserProfile(user.id)">查看主页</el-button>
									<el-button
										type="primary"
										:disabled="isFollowing(user.id)"
										:loading="followingUserId === user.id"
										@click="followUser(user.id)"
									>
										{{ isFollowing(user.id) ? '已关注' : '关注' }}
									</el-button>
								</div>
							</article>
						</div>

						<el-empty v-else description="暂无推荐朋友" :image-size="110">
							<template #description>
								<div class="empty-text">暂无推荐朋友，先去关注一些用户吧。</div>
							</template>
						</el-empty>
					</template>
				</el-card>
			</main>
		</div>
	</div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

import { postApi, userApi } from '../api'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const activeSection = ref('posts')
const postLoading = ref(false)
const friendLoading = ref(false)
const followingUserId = ref('')
const posts = ref([])
const suggestedUsers = ref([])
const currentPage = ref(1)
const numberOfPages = ref(1)
const pageSize = 2

const userId = computed(() => authStore.userId || '')
const isPostsSection = computed(() => activeSection.value === 'posts')
const sectionTitle = computed(() => (isPostsSection.value ? '推荐帖子' : '推荐朋友'))
const sectionDescription = computed(() =>
	isPostsSection.value
		? '根据你的关注关系实时更新的帖子推荐。'
		: '从你关注的人脉关系中继续发现值得关注的用户。',
)
const followingIds = computed(() => new Set((authStore.user?.following || []).map(String)))

function handleLogout() {
	authStore.logout()
	router.replace({ name: 'auth', params: { mode: 'login' } })
}

function authorText(name) {
	const value = String(name || '').trim()
	return value ? value.slice(0, 1).toUpperCase() : 'U'
}

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

async function fetchRecommendedPosts(page = 1) {
	if (!userId.value) {
		posts.value = []
		currentPage.value = 1
		numberOfPages.value = 1
		return
	}

	postLoading.value = true
	try {
		const response = await postApi.getAll({ id: userId.value, page })
		posts.value = normalizePosts(response?.posts || [])
		currentPage.value = response?.currentPage || page
		numberOfPages.value = response?.numberOfPages || 1
	} catch {
		posts.value = []
	} finally {
		postLoading.value = false
	}
}

async function fetchSuggestedFriends() {
	if (!userId.value) {
		suggestedUsers.value = []
		return
	}

	friendLoading.value = true
	try {
		const response = await userApi.getSuggested()
		suggestedUsers.value = normalizeUsers(response?.suggestedUsers || [])
	} catch {
		suggestedUsers.value = []
	} finally {
		friendLoading.value = false
	}
}

async function loadActiveSection() {
	if (isPostsSection.value) {
		await fetchRecommendedPosts(currentPage.value || 1)
		return
	}

	await fetchSuggestedFriends()
}

function handleSectionSelect(section) {
	activeSection.value = section
	if (section === 'posts') {
		currentPage.value = 1
	}
	loadActiveSection()
}

function handlePageChange(page) {
	if (page === currentPage.value || postLoading.value) {
		return
	}

	fetchRecommendedPosts(page)
}

function openPostDetail(postId) {
	if (!postId) {
		return
	}

	router.push({ name: 'post-detail', params: { id: postId } })
}

function openUserProfile(profileUserId) {
	if (!profileUserId) {
		return
	}

	router.push({ name: 'user-profile', params: { id: profileUserId } })
}

function isFollowing(profileUserId) {
	return followingIds.value.has(String(profileUserId || ''))
}

async function followUser(profileUserId) {
	if (!profileUserId || profileUserId === userId.value) {
		return
	}

	const wasFollowing = isFollowing(profileUserId)
	followingUserId.value = profileUserId
	try {
		const response = await userApi.follow(profileUserId)
		if (response?.firstUser) {
			authStore.updateUser(response.firstUser)
		}
		ElMessage.success(wasFollowing ? '已取消关注' : '关注成功')
		if (!wasFollowing && activeSection.value === 'friends') {
			await fetchSuggestedFriends()
		}
	} catch (error) {
		ElMessage.error(error?.response?.data?.error || '操作失败，请稍后重试')
	} finally {
		followingUserId.value = ''
	}
}

onMounted(() => {
	authStore.hydrate()
})

watch(
	userId,
	(nextUserId) => {
		if (!nextUserId) {
			posts.value = []
			suggestedUsers.value = []
			currentPage.value = 1
			numberOfPages.value = 1
			return
		}

		loadActiveSection()
	},
	{ immediate: true },
)
</script>

<style scoped>
.home-page {
	min-height: calc(100vh - 72px);
	padding: 20px;
	background:
		radial-gradient(circle at top right, rgba(37, 99, 235, 0.12), transparent 36%),
		radial-gradient(circle at bottom left, rgba(16, 185, 129, 0.12), transparent 32%),
		linear-gradient(180deg, #f7fbff 0%, #eef6ff 100%);
}

.home-layout {
	display: grid;
	grid-template-columns: 260px minmax(0, 1fr);
	gap: 20px;
	max-width: 1260px;
	margin: 0 auto;
}

.home-sidebar {
	display: flex;
	flex-direction: column;
	gap: 16px;
	padding: 20px;
	border-radius: 24px;
	background: rgba(255, 255, 255, 0.82);
	backdrop-filter: blur(18px);
	border: 1px solid rgba(148, 163, 184, 0.14);
	box-shadow: 0 18px 48px rgba(15, 23, 42, 0.08);
	position: sticky;
	top: 20px;
	height: fit-content;
}

.sidebar-brand {
	display: flex;
	align-items: center;
	gap: 12px;
	padding-bottom: 8px;
	border-bottom: 1px solid rgba(148, 163, 184, 0.16);
}

.brand-mark {
	width: 46px;
	height: 46px;
	border-radius: 14px;
	display: grid;
	place-items: center;
	background: linear-gradient(135deg, #2563eb, #10b981);
	color: #ffffff;
	font-weight: 800;
	letter-spacing: 0.04em;
}

.brand-name {
	font-size: 16px;
	font-weight: 800;
	color: #0f172a;
}

.brand-subtitle {
	margin-top: 2px;
	font-size: 12px;
	color: #64748b;
}

.section-menu {
	border-right: none;
	background: transparent;
}

.section-menu :deep(.el-menu-item) {
	border-radius: 14px;
	margin: 6px 0;
	height: 48px;
	font-weight: 600;
	color: #334155;
}

.section-menu :deep(.el-menu-item.is-active) {
	background: linear-gradient(135deg, rgba(37, 99, 235, 0.12), rgba(16, 185, 129, 0.12));
	color: #0f172a;
}

.sidebar-foot {
	padding-top: 8px;
	border-top: 1px solid rgba(148, 163, 184, 0.16);
}

.logout-button {
	width: 100%;
}

.home-main {
	min-width: 0;
}

.home-card {
	border-radius: 24px;
	border: 1px solid rgba(148, 163, 184, 0.14);
	box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08);
	background: rgba(255, 255, 255, 0.94);
}

.home-header {
	display: flex;
	align-items: flex-end;
	justify-content: space-between;
	gap: 16px;
}

.header-copy {
	min-width: 0;
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

.feed-skeleton {
	padding: 12px 2px 4px;
}

.feed-list,
.friend-list {
	display: flex;
	flex-direction: column;
	gap: 16px;
}

.feed-card,
.friend-card {
	padding: 18px;
	border-radius: 18px;
	background: #ffffff;
	border: 1px solid rgba(148, 163, 184, 0.16);
	box-shadow: 0 10px 24px rgba(15, 23, 42, 0.04);
	transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.feed-card:hover,
.friend-card:hover {
	transform: translateY(-1px);
	box-shadow: 0 14px 30px rgba(15, 23, 42, 0.08);
}

.feed-top,
.friend-top {
	display: flex;
	align-items: flex-start;
	justify-content: space-between;
	gap: 12px;
}

.feed-author,
.friend-author {
	display: flex;
	align-items: center;
	gap: 12px;
}

.feed-name,
.friend-name {
	font-size: 15px;
	font-weight: 700;
	color: #0f172a;
}

.friend-email,
.feed-meta {
	margin-top: 4px;
	font-size: 12px;
	color: #64748b;
}

.feed-counts,
.friend-stats {
	display: flex;
	gap: 8px;
	flex-wrap: wrap;
}

.friend-bio,
.post-message {
	margin: 12px 0 0;
	color: #334155;
	line-height: 1.8;
	white-space: pre-wrap;
}

.post-title {
	margin: 14px 0 8px;
	font-size: 18px;
	line-height: 1.35;
	color: #0f172a;
}

.post-image {
	display: block;
	width: 100%;
	max-height: 420px;
	margin-top: 14px;
	border-radius: 14px;
	overflow: hidden;
}

.feed-actions,
.friend-actions {
	margin-top: 16px;
	display: flex;
	justify-content: flex-end;
	gap: 10px;
	flex-wrap: wrap;
}

.pagination-wrap {
	margin-top: 18px;
	display: flex;
	justify-content: center;
}

.empty-text {
	color: #64748b;
}

@media (max-width: 960px) {
	.home-layout {
		grid-template-columns: 1fr;
	}

	.home-sidebar {
		position: static;
	}
}

@media (max-width: 640px) {
	.home-page {
		padding: 14px;
	}

	.home-header {
		flex-direction: column;
		align-items: flex-start;
	}

	.home-header h2 {
		font-size: 26px;
	}

	.feed-top,
	.friend-top {
		flex-direction: column;
	}

	.feed-actions,
	.friend-actions {
		justify-content: flex-start;
	}
}
</style>
