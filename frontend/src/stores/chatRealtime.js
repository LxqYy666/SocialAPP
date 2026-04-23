import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

function toWebSocketBaseUrl() {
	const customUrl = import.meta.env.VITE_CHAT_WS_BASE_URL
	if (customUrl) {
		return customUrl.replace(/\/$/, '')
	}

	const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
	try {
		const apiUrl = new URL(apiBaseUrl)
		const wsProtocol = apiUrl.protocol === 'https:' ? 'wss:' : 'ws:'
		return `${wsProtocol}//${apiUrl.hostname}:8081`
	} catch {
		return 'ws://localhost:8081'
	}
}

function normalizeIncomingMessage(raw) {
	if (!raw || typeof raw !== 'object') {
		return null
	}

	const sender = String(raw.sender || '').trim()
	const receiver = String(raw.receiver || '').trim()
	const content = String(raw.content || '').trim()
	if (!sender || !receiver || !content) {
		return null
	}

	return {
		id: raw.id || `rt-${Date.now()}-${Math.random().toString(16).slice(2)}`,
		sender,
		receiver,
		content,
	}
}

export const useChatRealtimeStore = defineStore('chat-realtime', () => {
	const socket = ref(null)
	const isConnected = ref(false)
	const currentUserId = ref('')
	const reconnectAttempts = ref(0)
	const reconnectTimer = ref(null)
	const shouldReconnect = ref(true)
	const incomingMessages = ref([])
	const onlineFriendIds = ref([])

	const hasIncomingMessage = computed(() => incomingMessages.value.length > 0)

	function clearReconnectTimer() {
		if (reconnectTimer.value) {
			window.clearTimeout(reconnectTimer.value)
			reconnectTimer.value = null
		}
	}

	function clearIncomingMessages() {
		incomingMessages.value = []
	}

	function setOnlineFriends(ids = []) {
		onlineFriendIds.value = Array.from(new Set(ids.map((id) => String(id)).filter(Boolean)))
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
				if (Array.isArray(payload?.onlineFriends)) {
					setOnlineFriends(payload.onlineFriends)
					return
				}

				const message = normalizeIncomingMessage(payload)
				if (message) {
					incomingMessages.value.push(message)
				}
			} catch {
				// ignore malformed payload
			}
		}

		ws.onclose = () => {
			isConnected.value = false
			socket.value = null
			setOnlineFriends([])

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
		setOnlineFriends([])
		clearIncomingMessages()

		if (socket.value) {
			socket.value.close()
			socket.value = null
		}
	}

	function sendMessage(receiver, content) {
		if (!receiver || !content || !socket.value || socket.value.readyState !== WebSocket.OPEN) {
			return false
		}

		socket.value.send(
			JSON.stringify({
				receiver,
				content,
			}),
		)
		return true
	}

	function consumeIncomingMessage() {
		if (!incomingMessages.value.length) {
			return null
		}
		return incomingMessages.value.shift() || null
	}

	return {
		isConnected,
		hasIncomingMessage,
		incomingMessages,
		onlineFriendIds,
		connect,
		disconnect,
		sendMessage,
		consumeIncomingMessage,
		clearIncomingMessages,
	}
})
