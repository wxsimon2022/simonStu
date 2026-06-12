<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const list = ref([])
const allPerms = ref([])
const loading = ref(false)
const createVisible = ref(false)
const editVisible = ref(false)
const form = ref({ name: '', description: '' })
const editForm = ref({ id: 0, description: '', status: 1, permission_ids: [] })

async function fetchRoles() {
  loading.value = true
  const res = await api('/api/role')
  if (res.code === 200) list.value = res.data.list
  loading.value = false
}

async function fetchPerms() {
  const res = await api('/api/permission')
  if (res.code === 200) {
    allPerms.value = res.data.list
  } else {
    console.error('权限列表加载失败', res)
  }
}

async function handleCreate() {
  if (!form.value.name) { ElMessage.warning('请输入角色标识'); return }
  const res = await api('/api/role', { method: 'POST', body: JSON.stringify(form.value) })
  if (res.code === 200) {
    ElMessage.success('创建成功')
    createVisible.value = false
    form.value = { name: '', description: '' }
    fetchRoles()
  } else ElMessage.error(res.message || '创建失败')
}

function openEdit(row) {
  editForm.value = {
    id: row.id,
    description: row.description,
    status: row.status,
    permission_ids: allPerms.value.filter(p => row.permissions?.includes(p.name)).map(p => p.id),
  }
  editVisible.value = true
}

async function handleEdit() {
  const res = await api('/api/role', { method: 'PUT', body: JSON.stringify(editForm.value) })
  if (res.code === 200) {
    ElMessage.success('更新成功')
    editVisible.value = false
    fetchRoles()
  } else ElMessage.error(res.message || '更新失败')
}

async function handleDelete(id) {
  try {
    await ElMessageBox.confirm('确定删除该角色？', '提示')
    const res = await api(`/api/role/${id}`, { method: 'DELETE' })
    if (res.code === 200) { ElMessage.success('已删除'); fetchRoles() }
  } catch (_) {}
}

onMounted(() => { fetchRoles(); fetchPerms() })
</script>

<template>
  <el-card>
    <template #header>
      <div class="flex items-center justify-between">
        <span>角色管理</span>
        <el-button type="primary" size="small" @click="createVisible = true">新建</el-button>
      </div>
    </template>
    <el-table :data="list" border stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="角色标识" />
      <el-table-column prop="description" label="描述" />
      <el-table-column prop="status" label="状态" width="70">
        <template #default="{ row }">
          <el-tag :type="row.status ? 'success' : 'danger'" size="small">{{ row.status ? '启用' : '禁用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="权限" min-width="200">
        <template #default="{ row }">
          <el-tag v-for="p in row.permissions" :key="p" size="small" style="margin:2px">{{ p }}</el-tag>
          <span v-if="!row.permissions?.length" style="color:#999;font-size:13px">无</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="130">
        <template #default="{ row }">
          <el-button type="primary" size="small" link @click="openEdit(row)">编辑</el-button>
          <el-button type="danger" size="small" link @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="createVisible" title="新建角色" width="420px">
    <el-form :model="form" label-width="80px">
      <el-form-item label="标识"><el-input v-model="form.name" placeholder="如 editor" /></el-form-item>
      <el-form-item label="描述"><el-input v-model="form.description" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="createVisible = false">取消</el-button>
      <el-button type="primary" @click="handleCreate">确定</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="editVisible" title="编辑角色" width="520px">
    <el-form :model="editForm" label-width="80px">
      <el-form-item label="描述"><el-input v-model="editForm.description" /></el-form-item>
      <el-form-item label="状态">
        <el-switch v-model="editForm.status" :active-value="1" :inactive-value="0" />
      </el-form-item>
      <el-form-item label="权限">
        <el-checkbox-group v-model="editForm.permission_ids">
          <el-checkbox
            v-for="p in allPerms"
            :key="p.id"
            :label="p.id"
            style="margin:4px 12px 4px 0"
          >
            {{ p.name }}
          </el-checkbox>
        </el-checkbox-group>
        <div v-if="!allPerms.length" style="color:#999;font-size:13px;margin-top:4px">暂无可用权限</div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="editVisible = false">取消</el-button>
      <el-button type="primary" @click="handleEdit">保存</el-button>
    </template>
  </el-dialog>
</template>
<VUE