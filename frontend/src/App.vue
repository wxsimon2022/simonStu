<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { Fold, Expand, Menu, Management } from '@element-plus/icons-vue'
import { api, getToken, removeToken } from './utils/request'

const router = useRouter()
const route = useRoute()
const isCollapse = ref(false)
const sidebarWidth = ref('220px')
const menus = ref([])

const isLoginPage = computed(() => route.path === '/login')

// 根据路由加载菜单
watch(() => route.path, async (path) => {
  if (path !== '/login' && getToken()) {
    const res = await api('/api/menus')
    if (res.code === 200) menus.value = res.data.menus || []
  }
}, { immediate: true })

// 根据数据库中配置的图标名称动态解析 Element Plus 图标组件
function getIcon(name) {
  return name && ElementPlusIconsVue[name] ? ElementPlusIconsVue[name] : Menu
}

async function handleLogout() {
  try { await api('/api/auth/logout', { method: 'POST' }) } catch (_) {}
  removeToken()
  router.push('/login')
}
</script>

<template>
  <router-view v-if="isLoginPage" />

  <el-container v-else style="height: 100vh">
    <el-aside :width="isCollapse ? '64px' : sidebarWidth">
      <div class="logo">
        <span v-if="!isCollapse">Gin Admin</span>
        <span v-else>GA</span>
      </div>
      <el-menu
        :router="true"
        :collapse="isCollapse"
        :default-active="route.path"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/">
          <el-icon><Management /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>

        <template v-for="item in menus" :key="item.id">
          <!-- 目录 → el-sub-menu（第一级） -->
          <el-sub-menu v-if="item.type === 'dir' && item.children?.length" :index="'dir-' + item.id">
            <template #title>
              <el-icon><component :is="getIcon(item.icon)" /></el-icon>
              <span>{{ item.description || item.name }}</span>
            </template>
            <!-- 菜单页 → el-menu-item（第二级） -->
            <el-menu-item v-for="child in item.children" :key="child.id" :index="child.route || ''">
              <template #title>
                <el-icon v-if="child.type === 'menu' && child.icon" :size="14"><component :is="getIcon(child.icon)" /></el-icon>
                <span>{{ child.description || child.name }}</span>
              </template>
            </el-menu-item>
          </el-sub-menu>
        </template>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="flex items-center justify-between px-4">
        <el-button @click="isCollapse = !isCollapse" text>
          <el-icon :size="20">
            <component :is="isCollapse ? Expand : Fold" />
          </el-icon>
        </el-button>
        <el-dropdown @command="handleLogout">
          <span class="cursor-pointer text-sm">管理员</span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style>
body { margin: 0; }
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  background: #2b3a4a;
}
.el-aside { background: #304156; overflow: hidden; transition: width 0.3s; }
.el-menu { border-right: none; }
.el-header {
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
}
.el-main { background: #f0f2f5; min-height: calc(100vh - 50px); }
.flex { display: flex; }
.items-center { align-items: center; }
.justify-between { justify-content: space-between; }
.px-4 { padding: 0 16px; }
.cursor-pointer { cursor: pointer; }
.text-sm { font-size: 14px; }
</style>
