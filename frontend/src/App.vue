<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Management, User, Setting, Fold, Expand } from '@element-plus/icons-vue'
import { api, removeToken } from './utils/request'

const router = useRouter()
const route = useRoute()
const isCollapse = ref(false)
const sidebarWidth = ref('220px')

const isLoginPage = computed(() => route.path === '/login')

async function handleLogout() {
  try { await api("/api/auth/logout", { method: "POST" }) } catch (_) {}
  removeToken()
  router.push('/login')
}
</script>

<template>
  <!-- 登录页：全屏居中，无侧边栏 -->
  <router-view v-if="isLoginPage" />

  <!-- 后台首页：侧边栏布局 -->
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
        <el-sub-menu index="system">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/admin">管理员管理</el-menu-item>
          <el-menu-item index="/role">角色管理</el-menu-item>
          <el-menu-item index="/permission">菜单列表</el-menu-item>
        </el-sub-menu>

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
