<template>
	<div class="post-detail-page">
		<el-card class="post-detail-card" shadow="never">
			<el-skeleton :loading="loading" animated>
				<template #default>
					<template v-if="post">
						<div class="detail-head">
							<div>
								<p class="eyebrow">帖子详情</p>
								<h1 class="post-title">{{ post.title }}</h1>
								<div class="post-meta">{{ formatPostDate(post.createdAt) }}</div>
							</div>
							<div class="meta-actions">
								<el-button :loading="liking" :type="isLiked ? 'success' : 'primary'" plain @click="toggleLike">
									{{ isLiked ? '已点赞' : '点赞' }}
								</el-button>
								<el-tag size="small" type="success" effect="plain">{{ likesCount }} 赞</el-tag>
								<el-tag size="small" type="info">{{ commentsCount }} 评论</el-tag>
							</div>
						</div>

						<el-divider />

						<div class="author-bar">
							<div class="author-info">
								<el-avatar :size="44">{{ authorText }}</el-avatar>
								<div>
									<div class="author-name">{{ post.name || '未命名作者' }}</div>
									<div class="author-id">ID: {{ post.creator }}</div>
								</div>
							</div>
							<el-button type="primary" plain @click="goCreatorProfile">查看作者主页</el-button>
						</div>

						<div class="post-body">
							<p class="post-message">{{ post.message }}</p>
							<el-image
								v-if="post.selectedFile"
								class="post-image"
								:src="post.selectedFile"
								fit="cover"
								:preview-src-list="[post.selectedFile]"
							/>
						</div>

						<el-divider />

						<div class="detail-foot">
							<div class="foot-item">
								<span class="label">帖子 ID</span>
								<span class="value">{{ post.id }}</span>
							</div>
							<div class="foot-item">
								<span class="label">点赞数量</span>
								<span class="value">{{ likesCount }}</span>
							</div>
							<div class="foot-item">
								<span class="label">评论数量</span>
								<span class="value">{{ commentsCount }}</span>
							</div>
						</div>

						<el-divider />

						<div class="comment-compose">
							<div class="section-head">
								<span class="section-title">发表评论</span>
								<el-tag size="small" type="info">发送后会刷新评论列表</el-tag>
							</div>

							<el-input
								v-model="commentDraft"
								type="textarea"
								:rows="4"
								maxlength="280"
								show-word-limit
								placeholder="写下你的评论"
							/>

							<div class="compose-actions">
								<el-button :disabled="submittingComment || !commentDraft.trim()" :loading="submittingComment" type="primary" @click="submitComment">
									发表评论
								</el-button>
							</div>
						</div>

						<div class="comments-section">
							<div class="section-head">
								<span class="section-title">评论列表</span>
								<el-tag size="small" type="info">{{ commentsCount }} 条</el-tag>
							</div>

							<div v-if="commentsList.length" class="comment-list">
								<div v-for="(comment, index) in commentsList" :key="`${index}-${comment.text}`" class="comment-item">
									<el-avatar :size="36">{{ commentText(comment.text) }}</el-avatar>
									<div class="comment-body">
										<div class="comment-meta">评论 {{ index + 1 }}</div>
										<p class="comment-text">{{ comment.text }}</p>
									</div>
								</div>
							</div>
							<el-empty v-else description="暂无评论" :image-size="90" />
						</div>

						<div class="detail-actions">
							<el-button @click="goBack">返回</el-button>
							<el-button type="primary" plain @click="goHome">回到首页</el-button>
						</div>
					</template>
					<el-empty v-else description="未找到帖子" :image-size="100" />
				</template>
			</el-skeleton>
		</el-card>
	</div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'

import { postApi } from '../api'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const loading = ref(false)
const liking = ref(false)
const submittingComment = ref(false)
const commentDraft = ref('')
const post = ref(null)

const postId = computed(() => String(route.params.id || ''))
const likesCount = computed(() => post.value?.likes?.length || 0)
const commentsCount = computed(() => post.value?.comments?.length || 0)
const currentUserId = computed(() => authStore.userId || '')
const isLiked = computed(() => Boolean(currentUserId.value && post.value?.likes?.includes(currentUserId.value)))
const commentsList = computed(() => normalizeComments(post.value?.comments || []))
const authorText = computed(() => (post.value?.name || 'U').trim().slice(0, 1).toUpperCase())

function normalizeComments(rawComments = []) {
	return rawComments.map((comment) => ({
		text: typeof comment === 'string' ? comment : comment?.value || comment?.text || '',
	}))
}

function commentText(value) {
	const text = String(value || '').trim()
	return text ? text.slice(0, 1).toUpperCase() : 'C'
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

async function loadPost() {
	if (!postId.value) {
		post.value = null
		return
	}

	loading.value = true
	try {
		const response = await postApi.getById(postId.value)
		post.value = response?.post || null
	} catch {
		post.value = null
	} finally {
		loading.value = false
	}
}

async function toggleLike() {
	if (!postId.value) {
		return
	}

	liking.value = true
	try {
		await postApi.like(postId.value)
		await loadPost()
	} catch (error) {
		ElMessage.error(error?.response?.data?.error || '点赞失败，请稍后重试')
	} finally {
		liking.value = false
	}
}

async function submitComment() {
	const value = commentDraft.value.trim()
	if (!value || !postId.value) {
		return
	}

	submittingComment.value = true
	try {
		await postApi.comment(postId.value, { value })
		commentDraft.value = ''
		await loadPost()
		ElMessage.success('评论已发布')
	} catch (error) {
		ElMessage.error(error?.response?.data?.error || '评论失败，请稍后重试')
	} finally {
		submittingComment.value = false
	}
}

function goCreatorProfile() {
	if (!post.value?.creator) {
		return
	}

	router.push({ name: 'user-profile', params: { id: post.value.creator } })
}

function goBack() {
	router.back()
}

function goHome() {
	router.push({ name: 'home' })
}

watch(postId, loadPost, { immediate: true })

onMounted(() => {
	authStore.hydrate()
	loadPost()
})
</script>

<style scoped>
.post-detail-page {
	min-height: calc(100vh - 72px);
	padding: 28px 20px;
	background:
		radial-gradient(circle at top right, rgba(37, 99, 235, 0.12), transparent 36%),
		radial-gradient(circle at bottom left, rgba(16, 185, 129, 0.12), transparent 32%),
		linear-gradient(180deg, #f7fbff 0%, #eef6ff 100%);
}

.post-detail-card {
	max-width: 920px;
	margin: 0 auto;
	padding: 6px;
	border-radius: 24px;
	border: 1px solid rgba(148, 163, 184, 0.14);
	box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08);
}

.detail-head {
	display: flex;
	align-items: flex-start;
	justify-content: space-between;
	gap: 16px;
	padding: 8px 6px 0;
}

.eyebrow {
	margin: 0 0 8px;
	font-size: 12px;
	font-weight: 700;
	letter-spacing: 0.14em;
	text-transform: uppercase;
	color: #64748b;
}

.post-title {
	margin: 0;
	font-size: clamp(1.5rem, 3vw, 2.1rem);
	line-height: 1.2;
	color: #0f172a;
}

.post-meta {
	margin-top: 10px;
	font-size: 13px;
	color: #64748b;
}

.meta-actions {
	display: flex;
	gap: 8px;
	flex-wrap: wrap;
	justify-content: flex-end;
}

.author-bar {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 16px;
}

.author-info {
	display: flex;
	align-items: center;
	gap: 12px;
}

.author-name {
	font-size: 16px;
	font-weight: 700;
	color: #0f172a;
}

.author-id {
	margin-top: 4px;
	font-size: 12px;
	color: #64748b;
}

.post-body {
	padding: 4px 0 0;
}

.post-message {
	margin: 0;
	color: #334155;
	line-height: 1.9;
	white-space: pre-wrap;
	font-size: 15px;
}

.post-image {
	display: block;
	width: 100%;
	max-height: 560px;
	margin-top: 16px;
	border-radius: 16px;
	overflow: hidden;
}

.detail-foot {
	display: grid;
	grid-template-columns: repeat(3, minmax(0, 1fr));
	gap: 12px;
}

.comments-section,
.comment-compose {
	margin-top: 4px;
}

.section-head {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 12px;
	margin-bottom: 12px;
}

.section-title {
	font-size: 12px;
	font-weight: 700;
	letter-spacing: 0.08em;
	text-transform: uppercase;
	color: #64748b;
}

.compose-actions {
	margin-top: 12px;
	display: flex;
	justify-content: flex-end;
}

.comment-list {
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.comment-item {
	display: flex;
	align-items: flex-start;
	gap: 12px;
	padding: 14px;
	border-radius: 14px;
	background: #f8fafc;
	border: 1px solid rgba(148, 163, 184, 0.18);
}

.comment-body {
	min-width: 0;
	flex: 1;
}

.comment-meta {
	font-size: 12px;
	font-weight: 700;
	color: #64748b;
	margin-bottom: 6px;
}

.comment-text {
	margin: 0;
	color: #334155;
	line-height: 1.8;
	white-space: pre-wrap;
}

.foot-item {
	padding: 14px;
	border-radius: 14px;
	background: #f8fafc;
	border: 1px solid rgba(148, 163, 184, 0.18);
}

.label {
	display: block;
	font-size: 12px;
	font-weight: 700;
	letter-spacing: 0.08em;
	text-transform: uppercase;
	color: #64748b;
}

.value {
	display: block;
	margin-top: 8px;
	font-size: 14px;
	font-weight: 600;
	color: #0f172a;
	word-break: break-all;
}

.detail-actions {
	margin-top: 20px;
	display: flex;
	justify-content: flex-end;
	gap: 10px;
}

@media (max-width: 760px) {
	.post-detail-page {
		padding: 16px 12px;
	}

	.detail-head,
	.author-bar {
		flex-direction: column;
	}

	.meta-actions,
	.detail-actions,
	.compose-actions {
		justify-content: flex-start;
	}

	.section-head {
		align-items: flex-start;
		flex-direction: column;
	}

	.detail-foot {
		grid-template-columns: 1fr;
	}
}
</style>
