import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { notificationApi } from '../api'

function inferNotificationType(details) {
	const text = String(details || '').toLowerCase()
	if (text.includes('following you')) {
		return 'follow'
	}
	if (text.includes('liked your post')) {
		return 'like'
	}
	if (text.includes('commented on your post')) {
		return 'comment'
	}
	return 'post'
}

function normalizeNotification(raw) {
	return {
		id: raw?.id || raw?._id || '',
		details: raw?.details || '',
		targetUserId: raw?.targetUserId || raw?.targetid || '',
		mainUserId: raw?.mainUserId || raw?.mainuid || '',
		type: raw?.type || inferNotificationType(raw?.details || ''),
		isReaded: typeof raw?.isReaded === 'boolean' ? raw.isReaded : Boolean(raw?.isreaded),
		createdAt: raw?.createdAt || '',
		notificationUser: {
			name: raw?.notificationUser?.name || raw?.user?.name || '',
			avatar: raw?.notificationUser?.avatar || raw?.user?.avatar || '',
		},
	}
}

function normalizeNotifications(rawNotifications = []) {
	return rawNotifications.map((item) => normalizeNotification(item))
}

function toWebSocketBaseUrl() {
	const customUrl = import.meta.env.VITE_NOTIFICATION_WS_BASE_URL
	if (customUrl) {
		return customUrl.replace(/\/$/, '')
	}

	const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
	try {
		const apiUrl = new URL(apiBaseUrl)
		const wsProtocol = apiUrl.protocol === 'https:' ? 'wss:' : 'ws:'
		return `${wsProtocol}//${apiUrl.hostname}:8082`
	} catch {
		return 'ws://localhost:8082'
	}
}

export const useNotificationRealtimeStore = defineStore('notification-realtime', () => {
	const notifications = ref([])
	const loading = ref(false)
	const isConnected = ref(false)
	const socket = ref(null)
	const reconnectTimer = ref(null)
	const reconnectAttempts = ref(0)
	const shouldReconnect = ref(true)
	const currentUserId = ref('')

	const unreadCount = computed(() => notifications.value.filter((item) => !item.isReaded).length)

	function clearReconnectTimer() {
		if (reconnectTimer.value) {
			window.clearTimeout(reconnectTimer.value)
			reconnectTimer.value = null
		}
	}

	function clear() {
		notifications.value = []
	}

	function upsertNotification(rawNotification, prepend = true) {
		const normalized = normalizeNotification(rawNotification)
		if (!normalized.id) {
			return
		}

		const index = notifications.value.findIndex((item) => item.id === normalized.id)
		if (index === -1) {
			if (prepend) {
				notifications.value.unshift(normalized)
				return
			}
			notifications.value.push(normalized)
			return
		}

		notifications.value[index] = {
			...notifications.value[index],
			...normalized,
		}
	}

	function connect(userId) {
		if (!userId) {
			return
		}

		if (currentUserId.value && currentUserId.value !== userId) {
			disconnect()
		}

		currentUserId.value = userId
		shouldReconnect.value = true
		clearReconnectTimer()

		if (socket.value && socket.value.readyState === WebSocket.OPEN) {
			return
		}

		const wsBaseUrl = toWebSocketBaseUrl()
		const wsUrl = `${wsBaseUrl}/ws/${encodeURIComponent(userId)}`
		const ws = new WebSocket(wsUrl)
		socket.value = ws

		ws.onopen = () => {
			isConnected.value = true
			reconnectAttempts.value = 0
		}

		ws.onmessage = (event) => {
			try {
				const payload = JSON.parse(event.data)
				upsertNotification(payload, true)
			} catch {
				// ignore malformed message payload
			}
		}

		ws.onclose = () => {
			isConnected.value = false
			socket.value = null

			if (!shouldReconnect.value || !currentUserId.value) {
				return
			}

			reconnectAttempts.value += 1
			const delay = Math.min(1000 * 2 ** (reconnectAttempts.value - 1), 10000)
			reconnectTimer.value = window.setTimeout(() => {
				connect(currentUserId.value)
			}, delay)
		}

		ws.onerror = () => {
			isConnected.value = false
		}
	}

	function disconnect() {
		shouldReconnect.value = false
		isConnected.value = false
		clearReconnectTimer()
		reconnectAttempts.value = 0
		currentUserId.value = ''

		if (socket.value) {
			socket.value.close()
			socket.value = null
		}
	}

	async function fetchNotifications() {
		if (!currentUserId.value) {
			notifications.value = []
			return
		}

		loading.value = true
		try {
			const response = await notificationApi.getAll()
			notifications.value = normalizeNotifications(response?.notifications || [])
		} catch {
			notifications.value = []
		} finally {
			loading.value = false
		}
	}

	async function markAllRead() {
		if (!unreadCount.value) {
			return false
		}

		await notificationApi.markRead()
		notifications.value = notifications.value.map((item) => ({
			...item,
			isReaded: true,
		}))
		return true
	}

	return {
		notifications,
		loading,
		isConnected,
		unreadCount,
		clear,
		connect,
		disconnect,
		fetchNotifications,
		markAllRead,
		upsertNotification,
	}
})
