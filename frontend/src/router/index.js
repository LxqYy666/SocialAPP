import { createRouter, createWebHistory } from 'vue-router'

import AuthView from '../views/AuthView.vue'
import HomeView from '../views/HomeView.vue'
import ChatView from '../views/ChatView.vue'
import SearchView from '../views/SearchView.vue'
import ProfileInfoView from '../views/ProfileInfoView.vue'
import ProfileRelationsView from '../views/ProfileRelationsView.vue'
import ProfileView from '../views/ProfileView.vue'
import ProfilePostsView from '../views/ProfilePostsView.vue'
import PostDetailView from '../views/PostDetailView.vue'
import UserProfileView from '../views/UserProfileView.vue'
import UserProfilePostsView from '../views/UserProfilePostsView.vue'
import UserProfileInfoView from '../views/UserProfileInfoView.vue'
import { useAuthStore } from '../stores/auth'
import { pinia } from '../stores/pinia'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/home',
    },
    {
      path: '/home',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true },
    },
    {
      path: '/messages',
      name: 'messages',
      component: ChatView,
      meta: { requiresAuth: true },
    },
    {
      path: '/search',
      name: 'search',
      component: SearchView,
      meta: { requiresAuth: true },
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true },
      redirect: { name: 'profile-info' },
      children: [
        {
          path: 'info',
          name: 'profile-info',
          component: ProfileInfoView,
          meta: { requiresAuth: true },
        },
        {
          path: 'posts',
          name: 'profile-posts',
          component: ProfilePostsView,
          meta: { requiresAuth: true },
        },
        {
          path: 'relations',
          name: 'profile-relations',
          component: ProfileRelationsView,
          meta: { requiresAuth: true },
        },
      ],
    },
    {
      path: '/auth/:mode(login|signup)?',
      name: 'auth',
      component: AuthView,
      meta: { guestOnly: true },
    },
    {
      path: '/users/:id',
      name: 'user-profile',
      component: UserProfileView,
      meta: { requiresAuth: true },
      redirect: { name: 'user-profile-info' },
      children: [
        {
          path: 'info',
          name: 'user-profile-info',
          component: UserProfileInfoView,
          meta: { requiresAuth: true },
        },
        {
          path: 'posts',
          name: 'user-profile-posts',
          component: UserProfilePostsView,
          meta: { requiresAuth: true },
        },
      ],
    },
    {
      path: '/posts/:id',
      name: 'post-detail',
      component: PostDetailView,
      meta: { requiresAuth: true },
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

router.beforeEach((to) => {
  const authStore = useAuthStore(pinia)

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return {
      name: 'auth',
      params: { mode: 'login' },
      query: { redirect: to.fullPath },
    }
  }

  if (to.meta.guestOnly && authStore.isAuthenticated) {
    return { name: 'home' }
  }

  return true
})

export default router
