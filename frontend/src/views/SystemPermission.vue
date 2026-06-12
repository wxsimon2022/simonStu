<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Folder, Menu, SwitchButton } from '@element-plus/icons-vue'

const treeData = ref([])
const loading = ref(false)
const selectedNode = ref(null)
const createVisible = ref(false)
const editVisible = ref(false)
const form = ref({ name: '', description: '', type: 'menu', parent_id: null })
const editForm = ref({ id: 0, name: '', description: '', type: 'menu', parent_id: null })

const treeProps = { children: 'children', label: 'name' }
const typeOptions = [
  { label: '目录', value: 'dir' },
  { label: '菜单', value: 'menu' },
  { label: '按钮', value: 'btn' },
]

async function fetchData() {
  loading.value = true
  const res = await api('/api/permission')
  if (res.code === 200) treeData.value = res.data.tree || []
  loading.value = false
}

function onNodeClick(data) {
  selectedNode.value = data
}

function openCreate() {
  form.value = { name: '', description: '', type: 'menu', parent_id: null }
  createVisible.value = true
}

function openEdit() {
  if (!selectedNode.value) { ElMessage.warning('请先选择一个菜单'); return }
  editForm.value = {
    id: selectedNode.value.id,
    name: selectedNode.value.name,
    description: selectedNode.value.description || '',
    type: selectedNode.value.type || 'menu',
    parent_id: selectedNode.value.parent_id ?? null,
  }
  editVisible.value = true
}

async function handleCreate() {
  if (!form.value.name) { ElMessage.warning('请输入菜单标识'); return }
  const body = { name: form.value.name, description: form.value.description || '', type: form.value.type || 'menu' }
  const res = await api('/api/permission', { method: 'POST', body: JSON.stringify(body) })
  if (res.code === 200) {
    ElMessage.success('创建成功'); createVisible.value = false; fetchData()
  } else ElMessage.error(res.message || '创建失败')
}

async function handleEdit() {
  const body = { id: editForm.value.id, name: editForm.value.name, description: editForm.value.description || '', type: editForm.value.type || 'menu', parent_id: editForm.value.parent_id ?? 0 }
  const res = await api('/api/permission', { method: 'PUT', body: JSON.stringify(body) })
  if (res.code === 200) {
    ElMessage.success('更新成功'); editVisible.value = false; selectedNode.value = null; fetchData()
  } else ElMessage.error(res.message || '更新失败')
}

async function handleDelete() {
  if (!selectedNode.value) { ElMessage.warning('请先选择一个菜单'); return }
  try {
    await ElMessageBox.confirm('确定删除菜单：' + selectedNode.value.name + '？', '提示')
    const res = await api(`/api/permission/${selectedNode.value.id}`, { method: 'DELETE' })
    if (res.code === 200) { ElMessage.success('已删除'); selectedNode.value = null; fetchData() }
  } catch (_) {}
}

onMounted(fetchData)
</script>

<template>
  <el-card>
    <template #header>
      <div class="flex items-center justify-between">
        <span>菜单列表</span>
        <div style="display:flex;gap:8px">
          <el-button type="primary" size="small" @click="openCreate">新建</el-button>
          <el-button size="small" :disabled="!selectedNode" @click="openEdit">编辑</el-button>
          <el-button type="danger" size="small" :disabled="!selectedNode" @click="handleDelete">删除</el-button>
        </div>
      </div>
    </template>
    <el-tree
      :data="treeData"
      :props="treeProps"
      node-key="id"
      default-expand-all
      highlight-current
      v-loading="loading"
      @node-click="onNodeClick"
    >
      <template #default="{ data }">
        <el-icon :size="16" style="margin-right:4px;vertical-align:middle">
          <Folder v-if="data.type === 'dir'" style="color:#e6a23c" />
          <Menu v-else-if="data.type === 'menu' || !data.type" style="color:#409eff" />
          <SwitchButton v-else style="color:#909399" />
        </el-icon>
        <span>{{ data.name }}</span>
        <span v-if="data.description" class="perm-desc">{{ data.description }}</span>
      </template>
    </el-tree>
    <el-empty v-if="!treeData.length && !loading" description="暂无数据" />
  </el-card>

  <el-dialog v-model="createVisible" title="新建菜单" width="420px">
    <el-form :model="form" label-width="80px">
      <el-form-item label="标识"><el-input v-model="form.name" placeholder="如 log:view" /></el-form-item>
      <el-form-item label="描述"><el-input v-model="form.description" /></el-form-item>
      <el-form-item label="类型">
        <el-select v-model="form.type" style="width:100%">
          <el-option v-for="t in typeOptions" :key="t.value" :label="t.label" :value="t.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="父菜单">
        <el-tree-select
          v-model="form.parent_id"
          :data="treeData"
          :props="{ children: 'children', label: 'name', value: 'id' }"
          placeholder="无（根节点）"
          clearable check-strictly filterable
          style="width:100%"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="createVisible = false">取消</el-button>
      <el-button type="primary" @click="handleCreate">确定</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="editVisible" title="编辑菜单" width="420px">
    <el-form :model="editForm" label-width="80px">
      <el-form-item label="标识"><el-input v-model="editForm.name" /></el-form-item>
      <el-form-item label="描述"><el-input v-model="editForm.description" /></el-form-item>
      <el-form-item label="类型">
        <el-select v-model="editForm.type" style="width:100%">
          <el-option v-for="t in typeOptions" :key="t.value" :label="t.label" :value="t.value" />
        </el-select>
      </el-form-item>
          <el-form-item label="父菜单">
            <el-tree-select
              v-model="editForm.parent_id"
              :data="treeData"
              :props="{ children: 'children', label: 'name', value: 'id' }"
              placeholder="无（根节点）"
              clearable check-strictly filterable
              style="width:100%"
            />
          </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="editVisible = false">取消</el-button>
      <el-button type="primary" @click="handleEdit">保存</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.perm-label { font-size: 14px; }
.perm-desc { color: #999; margin-left: 8px; font-size: 13px; }
</style>
