<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../utils/request'
import { Management, User, Collection, Menu } from '@element-plus/icons-vue'

const stats = ref(null)
const admins = ref([])
const loading = ref(true)

async function fetchStats() {
  try {
    const res = await api('/api/dashboard/stats')
    if (res.code === 200) stats.value = res.data
  } catch (_) {}
}

async function fetchAdmins() {
  try {
    const res = await api('/api/admin?page=1&page_size=6')
    if (res.code === 200) admins.value = res.data.list || []
  } catch (_) {}
}

onMounted(async () => {
  await Promise.all([fetchStats(), fetchAdmins()])
  loading.value = false
})
</script>

<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-body">
            <div class="stat-icon admin-icon"><el-icon :size="28"><Management /></el-icon></div>
            <div class="stat-info">
              <span class="stat-value">{{ stats?.admin_count ?? '-' }}</span>
              <span class="stat-label">管理员</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-body">
            <div class="stat-icon user-icon"><el-icon :size="28"><User /></el-icon></div>
            <div class="stat-info">
              <span class="stat-value">{{ stats?.user_count ?? '-' }}</span>
              <span class="stat-label">用户</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-body">
            <div class="stat-icon role-icon"><el-icon :size="28"><Collection /></el-icon></div>
            <div class="stat-info">
              <span class="stat-value">{{ stats?.role_count ?? '-' }}</span>
              <span class="stat-label">角色</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-body">
            <div class="stat-icon menu-icon"><el-icon :size="28"><Menu /></el-icon></div>
            <div class="stat-info">
              <span class="stat-value">{{ stats?.menu_count ?? '-' }}</span>
              <span class="stat-label">菜单</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近登录管理员列表 -->
    <el-card shadow="never" class="list-card" v-loading="loading">
      <template #header>
        <span class="card-title">管理员列表</span>
      </template>
      <el-table :data="admins" stripe style="width:100%" v-if="admins.length">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="real_name" label="姓名" />
        <el-table-column prop="create_time" label="创建时间" width="160" />
        <el-table-column label="状态" width="70">
          <template #default="{ row }">
            <el-tag :type="row.status ? 'success' : 'danger'" size="small">{{ row.status ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-else :image-size="60" description="暂无数据" />
    </el-card>
  </div>
</template>

<style scoped>
.dashboard {
  max-width: 1200px;
}
.stat-card {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  margin-bottom: 20px;
}
.stat-card :deep(.el-card__body) {
  padding: 20px;
}
.stat-body {
  display: flex;
  align-items: center;
  gap: 16px;
}
.stat-icon {
  width: 56px; height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.admin-icon { background: #e8f4ff; color: #409eff; }
.user-icon  { background: #e8f8e8; color: #67c23a; }
.role-icon  { background: #fef0e8; color: #e6a23c; }
.menu-icon  { background: #f4f0ff; color: #722ed1; }
.stat-info {
  display: flex;
  flex-direction: column;
}
.stat-value {
  font-size: 26px; font-weight: 700;
  color: #1d2129; line-height: 1.2;
  font-variant-numeric: tabular-nums;
}
.stat-label {
  font-size: 13px; color: #86909c; margin-top: 2px;
}
.list-card {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
.list-card :deep(.el-card__body) {
  padding: 0;
}
.list-card :deep(.el-card__header) {
  padding: 14px 24px;
  border-bottom: 1px solid #e5e6eb;
}
.card-title {
  font-size: 15px; font-weight: 600; color: #1d2129;
}
</style>
