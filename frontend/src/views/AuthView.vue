<template>
	<div class="auth-page">
		<div class="auth-backdrop"></div>
		<div class="auth-shell">
			<section class="auth-hero">
				<div class="hero-badge">SocialAPP</div>
				<h1>欢迎回来，继续连接你的社交空间</h1>
				<p>
					登录后即可查看动态、消息和通知，也可以快速创建新账号开始使用。
				</p>
				<ul>
					<li>统一登录与注册入口</li>
					<li>自动保存登录态</li>
					<li>支持路由直接跳转</li>
				</ul>
			</section>

			<section class="auth-card">
				<el-tabs v-model="activeTab" class="auth-tabs" @tab-change="handleTabChange">
					<el-tab-pane label="登录" name="login">
						<el-form
							ref="loginFormRef"
							:model="loginForm"
							:rules="loginRules"
							label-position="top"
							class="auth-form"
						>
							<el-form-item label="邮箱" prop="email">
								<el-input v-model="loginForm.email" type="email" size="large" placeholder="输入邮箱" />
							</el-form-item>

							<el-form-item label="密码" prop="password">
								<el-input
									v-model="loginForm.password"
									type="password"
									show-password
									size="large"
									placeholder="输入密码"
								/>
							</el-form-item>

							<el-alert v-if="authError" :title="authError" type="error" show-icon :closable="false" class="auth-alert" />

							<el-button type="primary" size="large" class="auth-submit" :loading="authStore.loading" @click="submitLogin">
								登录
							</el-button>
						</el-form>
					</el-tab-pane>

					<el-tab-pane label="注册" name="signup">
						<el-form
							ref="signupFormRef"
							:model="signupForm"
							:rules="signupRules"
							label-position="top"
							class="auth-form"
						>
							<el-row :gutter="12">
								<el-col :span="12">
									<el-form-item label="名" prop="firstName">
										<el-input v-model="signupForm.firstName" size="large" placeholder="First name" />
									</el-form-item>
								</el-col>
								<el-col :span="12">
									<el-form-item label="姓" prop="lastName">
										<el-input v-model="signupForm.lastName" size="large" placeholder="Last name" />
									</el-form-item>
								</el-col>
							</el-row>

							<el-form-item label="邮箱" prop="email">
								<el-input v-model="signupForm.email" type="email" size="large" placeholder="输入邮箱" />
							</el-form-item>

							<el-form-item label="密码" prop="password">
								<el-input
									v-model="signupForm.password"
									type="password"
									show-password
									size="large"
									placeholder="至少 5 位密码"
								/>
							</el-form-item>

							<el-alert v-if="authError" :title="authError" type="error" show-icon :closable="false" class="auth-alert" />

							<el-button type="primary" size="large" class="auth-submit" :loading="authStore.loading" @click="submitSignup">
								注册并登录
							</el-button>
						</el-form>
					</el-tab-pane>
				</el-tabs>
			</section>
		</div>
	</div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useAuthStore } from '../stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loginFormRef = ref()
const signupFormRef = ref()
const activeTab = ref(route.params.mode === 'signup' ? 'signup' : 'login')

const loginForm = reactive({
	email: '',
	password: '',
})

const signupForm = reactive({
	firstName: '',
	lastName: '',
	email: '',
	password: '',
})

const authError = computed(() => authStore.error)

const requiredRule = { required: true, trigger: 'blur' }

const loginRules = {
	email: [{ ...requiredRule, message: '请输入邮箱' }, { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }],
	password: [{ ...requiredRule, message: '请输入密码' }, { min: 5, message: '密码至少 5 位', trigger: 'blur' }],
}

const signupRules = {
	firstName: [{ ...requiredRule, message: '请输入名' }],
	lastName: [{ ...requiredRule, message: '请输入姓' }],
	email: [{ ...requiredRule, message: '请输入邮箱' }, { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }],
	password: [{ ...requiredRule, message: '请输入密码' }, { min: 5, message: '密码至少 5 位', trigger: 'blur' }],
}

function getRedirectTarget() {
	return typeof route.query.redirect === 'string' && route.query.redirect ? route.query.redirect : '/home'
}

async function submitLogin() {
	try {
		await loginFormRef.value?.validate()
		await authStore.signIn({ ...loginForm })
		await router.replace(getRedirectTarget())
	} catch {
		return
	}
}

async function submitSignup() {
	try {
		await signupFormRef.value?.validate()
		await authStore.signUp({ ...signupForm })
		await router.replace(getRedirectTarget())
	} catch {
		return
	}
}

function handleTabChange(tabName) {
	const nextMode = tabName === 'signup' ? 'signup' : 'login'
	if (route.params.mode !== nextMode) {
		router.replace({ name: 'auth', params: { mode: nextMode }, query: route.query })
	}
}

watch(
	() => route.params.mode,
	(mode) => {
		activeTab.value = mode === 'signup' ? 'signup' : 'login'
	},
)

onMounted(() => {
	authStore.hydrate()
	if (authStore.isAuthenticated) {
		router.replace('/home')
	}
})
</script>

<style scoped>
.auth-page {
	position: relative;
	min-height: 100vh;
	overflow: hidden;
	background:
		radial-gradient(circle at top left, rgba(96, 165, 250, 0.22), transparent 34%),
		radial-gradient(circle at bottom right, rgba(251, 191, 36, 0.16), transparent 30%),
		linear-gradient(135deg, #f8fbff 0%, #eef4ff 100%);
	padding: 40px 20px;
}

.auth-backdrop {
	position: absolute;
	inset: 0;
	background-image: linear-gradient(rgba(15, 23, 42, 0.03) 1px, transparent 1px), linear-gradient(90deg, rgba(15, 23, 42, 0.03) 1px, transparent 1px);
	background-size: 48px 48px;
	pointer-events: none;
}

.auth-shell {
	position: relative;
	z-index: 1;
	display: grid;
	grid-template-columns: minmax(0, 1.1fr) minmax(360px, 520px);
	gap: 28px;
	align-items: center;
	max-width: 1180px;
	margin: 0 auto;
	min-height: calc(100vh - 80px);
}

.auth-hero {
	color: #0f172a;
	padding: 24px;
}

.hero-badge {
	display: inline-flex;
	align-items: center;
	height: 36px;
	padding: 0 14px;
	margin-bottom: 18px;
	border-radius: 999px;
	background: rgba(255, 255, 255, 0.7);
	border: 1px solid rgba(148, 163, 184, 0.18);
	backdrop-filter: blur(10px);
	font-weight: 600;
	letter-spacing: 0.04em;
}

.auth-hero h1 {
	margin: 0;
	max-width: 520px;
	font-size: clamp(2.2rem, 5vw, 4.4rem);
	line-height: 1.02;
	letter-spacing: -0.04em;
}

.auth-hero p {
	max-width: 560px;
	margin: 18px 0 0;
	font-size: 1.05rem;
	line-height: 1.75;
	color: #475569;
}

.auth-hero ul {
	margin: 24px 0 0;
	padding: 0;
	list-style: none;
	color: #1e293b;
}

.auth-hero li {
	position: relative;
	padding-left: 22px;
	margin-top: 10px;
}

.auth-hero li::before {
	content: '';
	position: absolute;
	left: 0;
	top: 10px;
	width: 10px;
	height: 10px;
	border-radius: 50%;
	background: linear-gradient(135deg, #2563eb, #7c3aed);
}

.auth-card {
	padding: 24px;
	border-radius: 28px;
	background: rgba(255, 255, 255, 0.88);
	border: 1px solid rgba(148, 163, 184, 0.16);
	box-shadow: 0 24px 80px rgba(15, 23, 42, 0.12);
	backdrop-filter: blur(18px);
}

.auth-tabs :deep(.el-tabs__nav-wrap::after) {
	background-color: transparent;
}

.auth-tabs :deep(.el-tabs__item) {
	font-size: 15px;
	font-weight: 600;
}

.auth-form {
	padding-top: 12px;
}

.auth-alert {
	margin-bottom: 16px;
}

.auth-submit {
	width: 100%;
	height: 46px;
	border: 0;
	border-radius: 14px;
	background: linear-gradient(135deg, #2563eb 0%, #7c3aed 100%);
	box-shadow: 0 16px 30px rgba(37, 99, 235, 0.22);
}

.auth-submit:hover {
	filter: brightness(1.02);
}

@media (max-width: 960px) {
	.auth-shell {
		grid-template-columns: 1fr;
		min-height: auto;
	}

	.auth-hero {
		padding: 12px 4px 0;
	}
}

@media (max-width: 640px) {
	.auth-page {
		padding: 18px 12px;
	}

	.auth-card {
		padding: 18px 14px;
		border-radius: 22px;
	}
}
</style>