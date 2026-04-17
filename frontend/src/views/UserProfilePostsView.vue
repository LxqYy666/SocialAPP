<template>
	<div class="user-profile-posts-page">
		<div class="posts-panel">
			<div class="panel-head">
				<span class="label">发布内容</span>
				<el-tag size="small" type="info">{{ posts.length }} 篇</el-tag>
			</div>

			<div v-if="pagedPosts.length" class="pagination-head">
				<span class="pagination-tip">第 {{ currentPage }} / {{ totalPages }} 页</span>
			</div>

			<el-skeleton :loading="loading" animated>
				<template #default>
					<div v-if="pagedPosts.length" class="post-list">
						<article
							v-for="post in pagedPosts"
							:key="post.id"
							class="post-card"
							tabindex="0"
							role="button"
							@click="openPost(post.id)"
							@keydown.enter.prevent="openPost(post.id)"
							@keydown.space.prevent="openPost(post.id)"
						>
							<div class="post-top">
								<div>
									<h3 class="post-title">{{ post.title }}</h3>
									<div class="post-meta">{{ formatPostDate(post.createdAt) }}</div>
								</div>
								<el-tag size="small" type="success" effect="plain">{{ post.likes.length }} 赞</el-tag>
							</div>

							<p class="post-message">{{ post.message }}</p>

							<el-image
								v-if="post.selectedFile"
								class="post-image"
								:src="post.selectedFile"
								fit="cover"
								:preview-src-list="[post.selectedFile]"
								@click.stop
							/>
						</article>
					</div>
					<el-empty v-else description="暂无发布内容" :image-size="90" />
				</template>
			</el-skeleton>

			<div v-if="totalPages > 1" class="pagination-wrap">
				<el-pagination
					v-model:current-page="currentPage"
					background
					layout="prev, pager, next"
					:total="posts.length"
					:page-size="pageSize"
					@current-change="handlePageChange"
				/>
			</div>
		</div>
	</div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { userApi } from '../api'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const posts = ref([])
const currentPage = ref(1)
const pageSize = 6

const profileId = computed(() => String(route.params.id || ''))
const totalPages = computed(() => Math.max(1, Math.ceil(posts.value.length / pageSize)))
const pagedPosts = computed(() => {
	const start = (currentPage.value - 1) * pageSize
	return posts.value.slice(start, start + pageSize)
})

function normalizePosts(rawPosts = []) {
	return rawPosts.map((post) => ({
		id: post?.id || post?._id || '',
		title: post?.title || '未命名内容',
		message: post?.message || '',
		selectedFile: post?.selectedFile || '',
		likes: Array.isArray(post?.likes) ? post.likes : [],
		createdAt: post?.createdAt || '',
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

function openPost(postId) {
	if (!postId) {
		return
	}

	router.push({ name: 'post-detail', params: { id: postId } })
}

function handlePageChange(page) {
	currentPage.value = page
}

async function loadPosts() {
	if (!profileId.value) {
		posts.value = []
		currentPage.value = 1
		return
	}

	loading.value = true
	try {
		const response = await userApi.getById(profileId.value)
		posts.value = normalizePosts(response?.posts || [])
		currentPage.value = 1
	} catch {
		posts.value = []
	} finally {
		loading.value = false
	}
}

watch(profileId, loadPosts, { immediate: true })

onMounted(loadPosts)
</script>

<style scoped>
.user-profile-posts-page {
	min-height: 100%;
}

.posts-panel {
	padding: 18px;
	border-radius: 16px;
	background: #f8fafc;
	border: 1px solid rgba(148, 163, 184, 0.2);
}

.panel-head {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 12px;
}

.pagination-head {
	margin-bottom: 12px;
	display: flex;
	justify-content: flex-end;
}

.pagination-tip {
	font-size: 12px;
	color: #64748b;
}

.label {
	display: block;
	font-size: 12px;
	font-weight: 700;
	letter-spacing: 0.08em;
	text-transform: uppercase;
	color: #64748b;
}

.post-list {
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.post-card {
	padding: 16px;
	border-radius: 14px;
	background: #ffffff;
	border: 1px solid rgba(148, 163, 184, 0.16);
}

.post-top {
	display: flex;
	align-items: flex-start;
	justify-content: space-between;
	gap: 12px;
	margin-bottom: 10px;
}

.post-title {
	margin: 0;
	font-size: 16px;
	font-weight: 700;
	color: #0f172a;
}

.post-meta {
	margin-top: 6px;
	font-size: 12px;
	color: #64748b;
}

.post-message {
	margin: 0;
	color: #334155;
	line-height: 1.8;
	white-space: pre-wrap;
}

.post-image {
	display: block;
	width: 100%;
	max-height: 340px;
	margin-top: 12px;
	border-radius: 12px;
	overflow: hidden;
}

.pagination-wrap {
	margin-top: 16px;
	display: flex;
	justify-content: center;
}

@media (max-width: 760px) {
	.post-top {
		flex-direction: column;
	}
}
</style>