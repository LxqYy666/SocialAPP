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
				>
					<template #prefix>
						<svg viewBox="0 0 24 24" aria-hidden="true" class="search-icon">
							<path
								d="M10.5 4a6.5 6.5 0 1 1 0 13 6.5 6.5 0 0 1 0-13Zm0 2a4.5 4.5 0 1 0 0 9 4.5 4.5 0 0 0 0-9Zm5.97 8.56 3.97 3.97-1.41 1.41-3.97-3.97 1.41-1.41Z"
							/>
						</svg>
					</template>
				</el-input>
			</el-col>

			<el-col :xs="24" :sm="6" :md="7" class="nav-right">
				<el-space :size="10" alignment="center">
					<el-button class="icon-button" circle aria-label="消息">
						<svg viewBox="0 0 24 24" aria-hidden="true" class="action-icon">
							<path
								d="M4 4h16v12H7.17L4 19.17V4Zm2 2v8.34L6.34 14H18V6H6Zm2 2h8v2H8V8Zm0 4h6v2H8v-2Z"
							/>
						</svg>
					</el-button>

					<el-button class="icon-button" circle aria-label="通知">
						<svg viewBox="0 0 24 24" aria-hidden="true" class="action-icon">
							<path
								d="M12 3a5 5 0 0 0-5 5v2.28c0 .79-.25 1.56-.72 2.2L4 15v1h16v-1l-2.28-2.52c-.47-.64-.72-1.41-.72-2.2V8a5 5 0 0 0-5-5Zm0 18a2.5 2.5 0 0 0 2.45-2h-4.9A2.5 2.5 0 0 0 12 21Z"
							/>
						</svg>
					</el-button>

					<el-avatar class="avatar-button" :size="42" aria-label="个人头像">
						<svg viewBox="0 0 24 24" aria-hidden="true" class="action-icon">
							<path
								d="M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18Zm0 4a3.25 3.25 0 1 1 0 6.5A3.25 3.25 0 0 1 12 7Zm0 14a7.97 7.97 0 0 1-5.75-2.45A5.75 5.75 0 0 1 12 15.75a5.75 5.75 0 0 1 5.75 2.8A7.97 7.97 0 0 1 12 21Z"
							/>
						</svg>
					</el-avatar>
				</el-space>
			</el-col>
		</el-row>
	</header>
</template>

<script setup>
import { ref } from 'vue'

const searchText = ref('')
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

.avatar-button {
	width: 42px;
	height: 42px;
	padding: 0;
	background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
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
