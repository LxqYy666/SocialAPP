import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { authApi, clearAuthToken, getAuthToken, setAuthToken } from '../api'

const STORAGE_USER_KEY = 'user'

function readStoredUser() {
	const rawUser = localStorage.getItem(STORAGE_USER_KEY)
	if (!rawUser) {
		return null
	}

	try {
		return JSON.parse(rawUser)
	} catch {
		localStorage.removeItem(STORAGE_USER_KEY)
		return null
	}
}

function normalizeAuthResponse(response) {
	return {
		token: response?.token || '',
		user: response?.result || null,
	}
}

export const useAuthStore = defineStore('auth', () => {
	const token = ref(getAuthToken() || '')
	const user = ref(readStoredUser())
	const loading = ref(false)
	const error = ref('')

	const isAuthenticated = computed(() => Boolean(token.value))
	const userId = computed(() => user.value?._id || user.value?.id || '')

	function persistSession(nextToken, nextUser) {
		token.value = nextToken || ''
		user.value = nextUser || null

		if (token.value) {
			setAuthToken(token.value)
		} else {
			clearAuthToken()
		}

		if (user.value) {
			localStorage.setItem(STORAGE_USER_KEY, JSON.stringify(user.value))
		} else {
			localStorage.removeItem(STORAGE_USER_KEY)
		}
	}

	function setSessionFromResponse(response) {
		const session = normalizeAuthResponse(response)
		persistSession(session.token, session.user)
		return session
	}

	async function signIn(payload) {
		loading.value = true
		error.value = ''

		try {
			const response = await authApi.signIn(payload)
			return setSessionFromResponse(response)
		} catch (err) {
			error.value = err?.response?.data?.error || err?.response?.data?.message || '登录失败'
			throw err
		} finally {
			loading.value = false
		}
	}

	async function signUp(payload) {
		loading.value = true
		error.value = ''

		try {
			const response = await authApi.signUp(payload)
			return setSessionFromResponse(response)
		} catch (err) {
			error.value = err?.response?.data?.error || err?.response?.data?.message || '注册失败'
			throw err
		} finally {
			loading.value = false
		}
	}

	function logout() {
		persistSession('', null)
	}

	function updateUser(nextUser) {
		user.value = nextUser || null
		if (user.value) {
			localStorage.setItem(STORAGE_USER_KEY, JSON.stringify(user.value))
		} else {
			localStorage.removeItem(STORAGE_USER_KEY)
		}
	}

	function hydrate() {
		const storedToken = getAuthToken()
		const storedUser = readStoredUser()
		persistSession(storedToken || '', storedUser)
	}

	return {
		token,
		user,
		loading,
		error,
		isAuthenticated,
		userId,
		hydrate,
		signIn,
		signUp,
		logout,
		updateUser,
		persistSession,
		setSessionFromResponse,
	}
})