<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../utils/request'

const users = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const loading = ref(false)

async function fetchUsers() {
  loading.value = true
  try {
    const res = await api(`/api/user?page=${page.value}&size=${size.value}`)
    if (res.code === 200) {
      users.value = res.data.list
      total.value = res.data.total
    }
  } catch (e) {
    console.error('获取用户列表失败', e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchUsers)
</script>

<template>
  <div>
    <el-card>
      <template #header><span>用户管理</span></template>

      <el-table :data="users" border stripe v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="is_admin" label="管理员" width="90">
          <template #default="{ row }">
            <el-tag :type="row.is_admin ? 'success' : 'info'" size="small">
              {{ row.is_admin ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="120" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button type="primary" size="small" link>编辑</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="size"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        background
        @current-change="fetchUsers"
        @size-change="fetchUsers"
        style="margin-top: 20px; justify-content: flex-end; display: flex"
      />
    </el-card>
  </div>
</template>
