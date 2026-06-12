<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const list = ref([])
const allRoles = ref([])
const loading = ref(false)
const createVisible = ref(false)
const editVisible = ref(false)
const form = ref({ username: '', password: '', real_name: '', phone: '', email: '', role_ids: [] })
const editForm = ref({ id: 0, real_name: '', phone: '', email: '', password: '', role_ids: [] })

async function fetchData() {
  loading.value = true
  const res = await api('/api/admin')
  if (res.code === 200) list.value = res.data.list
  loading.value = false
}

async function fetchRoles() {
  const res = await api('/api/role')
  if (res.code === 200) allRoles.value = res.data.list
}

async function handleCreate() {
  const body = {
    username: form.value.username,
    password: form.value.password,
    real_name: form.value.real_name || '',
    phone: form.value.phone || '',
    email: form.value.email || '',
    role_ids: form.value.role_ids || [],
  }
  const res = await api('/api/admin', { method: 'POST', body: JSON.stringify(body) })
  if (res.code === 200) {
    ElMessage.success('创建成功')
    createVisible.value = false
    form.value = { username: '', password: '', real_name: '', phone: '', email: '', role_ids: [] }
    fetchData()
  } else ElMessage.error(res.message || '创建失败')
}

function openEdit(row) {
  editForm.value = {
    id: row.id,
    real_name: row.real_name || '',
    phone: row.phone || '',
    email: row.email || '',
    password: '',
    role_ids: row.role_ids || [],
  }
  editVisible.value = true
}

async function handleEdit() {
  const body = {
    id: editForm.value.id,
    real_name: editForm.value.real_name,
    phone: editForm.value.phone,
    email: editForm.value.email,
    role_ids: editForm.value.role_ids || [],
  }
  if (editForm.value.password) body.password = editForm.value.password
  const res = await api('/api/admin', { method: 'PUT', body: JSON.stringify(body) })
  if (res.code === 200) {
    ElMessage.success('更新成功')
    editVisible.value = false
    fetchData()
  } else ElMessage.error(res.message || '更新失败')
}

async function handleDelete(id) {
  try {
    await ElMessageBox.confirm('确定删除该管理员？', '提示')
    const res = await api(`/api/admin/${id}`, { method: 'DELETE' })
    if (res.code === 200) { ElMessage.success('已删除'); fetchData() }
  } catch (_) {}
}

onMounted(() => { fetchData(); fetchRoles() })
</script>

<template>
  <el-card>
    <template #header>
      <div class="flex items-center justify-between">
        <span>管理员管理</span>
        <el-button type="primary" size="small" @click="createVisible = true">新建</el-button>
      </div>
    </template>
    <el-table :data="list" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="real_name" label="姓名" />
      <el-table-column label="角色" min-width="140">
        <template #default="{ row }">
          <el-tag v-for="r in row.roles" :key="r" size="small" style="margin: 2px">{{ r }}</el-tag>
          <span v-if="!row.roles?.length" style="color:#999;font-size:13px">无</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="70">
        <template #default="{ row }">
          <el-tag :type="row.status ? 'success' : 'danger'" size="small">{{ row.status ? '启用' : '禁用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="create_time" label="创建时间" width="120" />
      <el-table-column label="操作" width="130">
        <template #default="{ row }">
          <el-button type="primary" size="small" link @click="openEdit(row)">编辑</el-button>
          <el-button type="danger" size="small" link @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="createVisible" title="新建管理员" width="480px">
    <el-form :model="form" label-width="80px">
      <el-form-item label="用户名"><el-input v-model="form.username" /></el-form-item>
      <el-form-item label="密码"><el-input v-model="form.password" type="password" /></el-form-item>
      <el-form-item label="姓名"><el-input v-model="form.real_name" /></el-form-item>
      <el-form-item label="手机"><el-input v-model="form.phone" /></el-form-item>
      <el-form-item label="邮箱"><el-input v-model="form.email" /></el-form-item>
      <el-form-item label="角色">
        <el-checkbox-group v-model="form.role_ids">
          <el-checkbox v-for="r in allRoles" :key="r.id" :label="r.id" style="margin:4px 12px 4px 0">{{ r.name }}</el-checkbox>
        </el-checkbox-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="createVisible = false">取消</el-button>
      <el-button type="primary" @click="handleCreate">确定</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="editVisible" title="编辑管理员" width="480px">
    <el-form :model="editForm" label-width="80px">
      <el-form-item label="姓名"><el-input v-model="editForm.real_name" /></el-form-item>
      <el-form-item label="手机"><el-input v-model="editForm.phone" /></el-form-item>
      <el-form-item label="邮箱"><el-input v-model="editForm.email" /></el-form-item>
      <el-form-item label="新密码"><el-input v-model="editForm.password" type="password" placeholder="留空不修改" /></el-form-item>
      <el-form-item label="角色">
        <el-checkbox-group v-model="editForm.role_ids">
          <el-checkbox v-for="r in allRoles" :key="r.id" :label="r.id" style="margin:4px 12px 4px 0">{{ r.name }}</el-checkbox>
        </el-checkbox-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="editVisible = false">取消</el-button>
      <el-button type="primary" @click="handleEdit">保存</el-button>
    </template>
  </el-dialog>
</template>
<VUE