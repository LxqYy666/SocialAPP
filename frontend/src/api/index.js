import request from '../utils/request'

export const authApi = {
	signUp(data) {
		return request.post('/user/signup', data)
	},
	signIn(data) {
		return request.post('/user/signin', data)
	},
}

export const userApi = {
	getById(id) {
		return request.get(`/user/getuser/${id}`)
	},
	update(id, data) {
		return request.patch(`/user/update/${id}`, data)
	},
	follow(id) {
		return request.patch(`/user/${id}/following`)
	},
	getSuggested() {
		return request.get('/user/sug/user')
	},
	delete(id) {
		return request.delete(`/user/delete/${id}`)
	},
}

export const postApi = {
	create(data) {
		return request.post('/post/create', data)
	},
	getById(id) {
		return request.get(`/post/get/${id}`)
	},
	getAll({ id, page = 1 }) {
		return request.get('/post/all', { params: { id, page } })
	},
	search(searchQuery) {
		return request.get('/post/search', { params: { searchQuery } })
	},
	update(id, data) {
		return request.patch(`/post/${id}`, data)
	},
	comment(id, data) {
		return request.post(`/post/${id}/commentPost`, data)
	},
	like(id) {
		return request.post(`/post/${id}/likePost`)
	},
	delete(id) {
		return request.delete(`/post/${id}`)
	},
}

export const chatApi = {
	sendMessage(data) {
		return request.post('/chat/sendMsg', data)
	},
	getMessages({ from, firstuid, seconduid }) {
		return request.get('/chat/getMsgByNums', { params: { from, firstuid, seconduid } })
	},
	getUnreadMessages() {
		return request.get('/chat/getUserUnReadedMsg')
	},
	markUnreadMessages(otherUserId) {
		return request.patch('/chat/maskUnReadedMsg', null, { params: { otherUserId } })
	},
}

export const notificationApi = {
	getAll() {
		return request.get('/notification/get')
	},
	markRead() {
		return request.patch('/notification/markread')
	},
}

export function setAuthToken(token) {
	if (token) {
		localStorage.setItem('token', token)
		return
	}
	localStorage.removeItem('token')
}

export function getAuthToken() {
	return localStorage.getItem('token')
}

export function clearAuthToken() {
	localStorage.removeItem('token')
}