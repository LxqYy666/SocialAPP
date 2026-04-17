<template>
	<div class="profile-info-page">
		<div class="profile-grid">
			<div class="profile-item detail-list">
				<div class="detail-row">
					<span class="label">用户名称</span>
					<span class="detail-value">{{ displayName }}</span>
				</div>
				<div class="detail-row">
					<span class="label">邮箱地址</span>
					<span class="detail-value">{{ displayEmail }}</span>
				</div>
				<div class="detail-row">
					<span class="label">用户 ID</span>
					<span class="detail-value">{{ displayUserId }}</span>
				</div>
				<div class="detail-row">
					<span class="label">个人简介</span>
					<p class="detail-bio">{{ displayBio }}</p>
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

		<el-divider />

		<el-form ref="formRef" :model="formModel" :rules="formRules" label-position="top" class="edit-form">
			<el-row :gutter="16">
				<el-col :xs="24" :md="12">
					<el-form-item label="用户名称" prop="name">
						<el-input v-model="formModel.name" placeholder="请输入用户名称" clearable />
					</el-form-item>
				</el-col>
				<el-col :xs="24" :md="12">
					<el-form-item label="邮箱地址">
						<el-input :model-value="displayEmail" disabled />
					</el-form-item>
				</el-col>
			</el-row>

			<el-form-item label="头像地址" prop="imageUrl">
				<el-input
					v-model="formModel.imageUrl"
					placeholder="请输入图片 URL，例如 https://example.com/avatar.png"
					clearable
				/>
			</el-form-item>

			<el-form-item label="个人简介" prop="bio">
				<el-input
					v-model="formModel.bio"
					type="textarea"
					:rows="4"
					maxlength="280"
					show-word-limit
					placeholder="介绍一下自己吧"
				/>
			</el-form-item>
		</el-form>

		<div class="profile-actions">
			<el-button :disabled="saving" @click="resetForm">重置</el-button>
			<el-button type="primary" :loading="saving" @click="saveProfile">保存修改</el-button>
		</div>
	</div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'

import { userApi } from '../api'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const formRef = ref()
const saving = ref(false)

const user = computed(() => authStore.user || {})
const displayName = computed(() => user.value?.name || '未命名用户')
const displayEmail = computed(() => user.value?.email || '暂无邮箱信息')
const displayBio = computed(() => user.value?.bio || '这个人很神秘，还没有填写简介。')
const displayUserId = computed(() => user.value?.id || user.value?._id || '-')
const followersCount = computed(() => user.value?.followers?.length || 0)
const followingCount = computed(() => user.value?.following?.length || 0)

const formModel = reactive({
	name: '',
	bio: '',
	imageUrl: '',
})

const formRules = {
	name: [{ required: true, message: '请输入用户名称', trigger: 'blur' }],
	imageUrl: [
		{
			validator: (_, value, callback) => {
				if (!value) {
					callback()
					return
				}
				const isValid = /^https?:\/\//.test(value)
				callback(isValid ? undefined : new Error('头像地址需以 http:// 或 https:// 开头'))
			},
			trigger: 'blur',
		},
	],
}

function syncFormFromUser() {
	formModel.name = user.value?.name || ''
	formModel.bio = user.value?.bio || ''
	formModel.imageUrl = user.value?.imageUrl || ''
}

async function refreshUser() {
	if (!authStore.userId) {
		return
	}

	try {
		const response = await userApi.getById(authStore.userId)
		if (response?.user) {
			authStore.updateUser(response.user)
			syncFormFromUser()
		}
	} catch {
		// Ignore refresh error to avoid blocking form usage with cached user data.
	}
}

function resetForm() {
	syncFormFromUser()
	formRef.value?.clearValidate()
}

async function saveProfile() {
	if (!authStore.userId) {
		ElMessage.error('未获取到用户信息，请重新登录后再试')
		return
	}

	try {
		await formRef.value?.validate()
	} catch {
		return
	}

	saving.value = true
	try {
		const payload = {
			name: formModel.name.trim(),
			bio: formModel.bio.trim(),
			imageUrl: formModel.imageUrl.trim(),
		}
		const response = await userApi.update(authStore.userId, payload)
		const updated = response?.data
		if (updated) {
			authStore.updateUser(updated)
			syncFormFromUser()
		}
		ElMessage.success('个人信息已更新')
	} catch (err) {
		ElMessage.error(err?.response?.data?.error || '保存失败，请稍后重试')
	} finally {
		saving.value = false
	}
}

onMounted(async () => {
	authStore.hydrate()
	syncFormFromUser()
	await refreshUser()
})
</script>

<style scoped>
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

.edit-form {
	margin-top: 6px;
}

.profile-actions {
	margin-top: 20px;
	display: flex;
	justify-content: flex-end;
	gap: 10px;
}

@media (max-width: 760px) {
	.profile-grid {
		grid-template-columns: 1fr;
	}
}
</style>