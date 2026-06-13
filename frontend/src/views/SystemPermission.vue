<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { api } from '../utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { Folder, FolderOpened, Menu, SwitchButton, Search } from '@element-plus/icons-vue'

const treeRef = ref(null)
const treeData = ref([])
const loading = ref(false)
const selectedNode = ref(null)
const filterText = ref('')
const createVisible = ref(false)
const editVisible = ref(false)
const contextMenu = ref({ visible: false, x: 0, y: 0 })
const form = ref({ name: '', description: '', type: 'menu', icon: '', parent_id: null })
const editForm = ref({ id: 0, name: '', description: '', type: 'menu', icon: '', parent_id: null })

const treeProps = { children: 'children', label: 'name' }
// 常用 Element Plus 图标列表
const iconOptions = [
  '', 'User', 'UserFilled', 'Management', 'Setting', 'Menu',
  'Folder', 'FolderOpened', 'HomeFilled', 'Tools',
  'Histogram', 'DataAnalysis', 'Aim', 'Bell',
  'ChatDotSquare', 'CircleCheck', 'Document', 'Notebook',
  'Files', 'Search', 'Plus', 'Edit', 'Delete',
  'SwitchButton', 'Refresh', 'Share', 'Upload',
  'Flag', 'Star', 'Collection', 'List',
]

const typeOptions = [
  { label: '目录', value: 'dir' },
  { label: '菜单', value: 'menu' },
  { label: '按钮', value: 'btn' },
]

function typeLabel(t) { return { dir: '目录', menu: '菜单', btn: '按钮' }[t] || '菜单' }
function typeTagType(t) { return { dir: 'warning', menu: 'primary', btn: 'info' }[t] || 'primary' }

// 树节点图标：优先用数据库中配置的 icon，否则按 type 显示默认图标
function getNodeIcon(data) {
  if (data.icon && ElementPlusIconsVue[data.icon]) return ElementPlusIconsVue[data.icon]
  const fallback = { dir: Folder, menu: Menu, btn: SwitchButton }
  return fallback[data.type] || Menu
}

function countNodes(nodes) {
  let c = 0
  for (const n of nodes) { c++; if (n.children?.length) c += countNodes(n.children) }
  return c
}
const flatCount = computed(() => countNodes(treeData.value))

watch(filterText, (v) => { treeRef.value?.filter(v) })
function filterNode(value, data) {
  if (!value) return true
  const q = value.toLowerCase()
  return data.name.toLowerCase().includes(q) || (data.description || '').toLowerCase().includes(q)
}

function setExpand(expand) {
  treeRef.value?.store?.nodesMap?.forEach(node => {
    if (expand) node.expand()
    else node.collapse()
  })
}

function onKeyDown(e) {
  if (e.key === 'Enter' && selectedNode.value && !createVisible.value && !editVisible.value) openEdit()
  if ((e.key === 'Delete' || e.key === 'Backspace') && selectedNode.value && !createVisible.value && !editVisible.value) handleDelete()
}
function closeCtx() { contextMenu.value.visible = false }

onMounted(() => {
  document.addEventListener('click', closeCtx)
  document.addEventListener('keydown', onKeyDown)
  fetchData()
})
onUnmounted(() => {
  document.removeEventListener('click', closeCtx)
  document.removeEventListener('keydown', onKeyDown)
})

async function fetchData() {
  loading.value = true
  const res = await api('/api/permission')
  if (res.code === 200) { treeData.value = res.data.tree || [] } else { console.error("菜单列表加载失败", res) }
  loading.value = false
}

function onNodeClick(data) { selectedNode.value = data }
function onContextMenu(e, data) {
  e.preventDefault()
  selectedNode.value = data
  contextMenu.value = { visible: true, x: e.clientX, y: e.clientY }
}

function openCreate() {
  form.value = { name: '', description: '', type: 'menu', icon: '', parent_id: null }
  createVisible.value = true
}
function createChild() {
  contextMenu.value.visible = false
  const pid = selectedNode.value?.id || null
  form.value = { name: '', description: '', type: 'menu', icon: '', parent_id: pid }
  createVisible.value = true
}
function openEdit() {
  contextMenu.value.visible = false
  if (!selectedNode.value) { ElMessage.warning('请先选择一个菜单'); return }
  editForm.value = {
    id: selectedNode.value.id, name: selectedNode.value.name,
    description: selectedNode.value.description || '', type: selectedNode.value.type || 'menu',
    icon: selectedNode.value.icon || '',
    parent_id: selectedNode.value.parent_id ?? null,
  }
  editVisible.value = true
}

async function handleCreate() {
  if (!form.value.name) { ElMessage.warning('请输入菜单标识'); return }
  const body = { name: form.value.name, description: form.value.description || '', type: form.value.type || 'menu', icon: form.value.icon || '' }
  const res = await api('/api/permission', { method: 'POST', body: JSON.stringify(body) })
  if (res.code === 200) { ElMessage.success('创建成功'); createVisible.value = false; fetchData() }
  else ElMessage.error(res.message || '创建失败')
}

async function handleEdit() {
  const body = { id: editForm.value.id, name: editForm.value.name, description: editForm.value.description || '', type: editForm.value.type || 'menu', icon: editForm.value.icon || '', parent_id: editForm.value.parent_id ?? 0 }
  const res = await api('/api/permission', { method: 'PUT', body: JSON.stringify(body) })
  if (res.code === 200) { ElMessage.success('更新成功'); editVisible.value = false; selectedNode.value = null; fetchData() }
  else ElMessage.error(res.message || '更新失败')
}

async function handleDelete() {
  contextMenu.value.visible = false
  if (!selectedNode.value) { ElMessage.warning('请先选择一个菜单'); return }
  try {
    await ElMessageBox.confirm('确定删除菜单：' + selectedNode.value.name + '？', '提示')
    const res = await api(`/api/permission/${selectedNode.value.id}`, { method: 'DELETE' })
    if (res.code === 200) { ElMessage.success('已删除'); selectedNode.value = null; fetchData() }
  } catch (_) {}
}
</script>

<template>
  <el-card shadow="never" class="page-card">
    <template #header>
      <div class="card-header">
        <span class="card-title">菜单列表</span>
        <div class="header-actions">
          <el-button type="primary" size="small" @click="openCreate">新建</el-button>
          <el-button size="small" :disabled="!selectedNode" @click="openEdit">编辑</el-button>
          <el-button type="danger" size="small" :disabled="!selectedNode" @click="handleDelete">删除</el-button>
        </div>
      </div>
    </template>

    <div class="toolbar">
      <el-input v-model="filterText" placeholder="搜索菜单名称或描述…" clearable :prefix-icon="Search" size="small" style="width:260px" />
      <div class="toolbar-right">
        <el-tooltip content="展开全部">
          <el-button size="small" class="toolbar-btn" @click="setExpand(true)"><el-icon><FolderOpened /></el-icon></el-button>
        </el-tooltip>
        <el-tooltip content="折叠全部">
          <el-button size="small" class="toolbar-btn" @click="setExpand(false)"><el-icon><Folder /></el-icon></el-button>
        </el-tooltip>
        <span class="node-count">共 {{ flatCount }} 项</span>
      </div>
    </div>

    <div class="type-legend">
      <span class="legend-item"><el-icon :size="14"><Folder style="color:#e6a23c" /></el-icon> 目录</span>
      <span class="legend-item"><el-icon :size="14"><Menu style="color:#409eff" /></el-icon> 菜单</span>
      <span class="legend-item"><el-icon :size="14"><SwitchButton style="color:#909399" /></el-icon> 按钮</span>
      <span class="legend-hint">右键更多操作 · 回车编辑 · Delete 删除</span>
    </div>

    <el-tree
      ref="treeRef"
      :data="treeData"
      :props="treeProps"
      :filter-node-method="filterNode"
      node-key="id"
      default-expand-all
      highlight-current
      v-loading="loading"
      @node-click="onNodeClick"
      @node-contextmenu="onContextMenu"
      class="menu-tree"
    >
      <template #default="{ data }">
        <div class="tree-node">
          <div class="node-info">
            <el-icon :size="16" class="node-icon">
              <component :is="getNodeIcon(data)" />
            </el-icon>
            <span class="node-name">{{ data.name }}</span>
            <span v-if="data.description" class="node-desc">{{ data.description }}</span>
            <span v-if="data.icon" class="node-icon-name">{{ data.icon }}</span>
          </div>
          <el-tag :type="typeTagType(data.type)" size="small" effect="plain" class="node-tag">
            {{ typeLabel(data.type) }}
          </el-tag>
        </div>
      </template>
    </el-tree>

    <el-empty v-if="!treeData.length && !loading" :image-size="80" description="暂无菜单数据，点击上方「新建」创建" />
  </el-card>

  <teleport to="body">
    <transition name="ctx-fade">
      <div v-if="contextMenu.visible" class="ctx-menu" :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }" @click.stop>
        <div class="ctx-item" @click="createChild">
          <span class="ctx-icon">+</span> 新建子菜单
        </div>
        <div class="ctx-divider"></div>
        <div class="ctx-item" @click="openEdit">
          <span class="ctx-icon">✎</span> 编辑
        </div>
        <div class="ctx-divider"></div>
        <div class="ctx-item ctx-danger" @click="handleDelete">
          <span class="ctx-icon">✕</span> 删除
        </div>
      </div>
    </transition>
  </teleport>

  <el-dialog v-model="createVisible" title="新建菜单" width="420px" :close-on-click-modal="false">
    <el-form :model="form" label-width="80px">
      <el-form-item label="标识"><el-input v-model="form.name" placeholder="如 log:view" /></el-form-item>
      <el-form-item label="描述"><el-input v-model="form.description" /></el-form-item>
      <el-form-item label="类型"><el-select v-model="form.type" style="width:100%"><el-option v-for="t in typeOptions" :key="t.value" :label="t.label" :value="t.value" /></el-select></el-form-item>
      <el-form-item label="图标">
        <el-select v-model="form.icon" filterable clearable placeholder="选择图标（仅目录/菜单）" style="width:100%">
          <el-option v-for="ico in iconOptions" :key="ico" :label="ico || '无'" :value="ico">
            <span style="display:flex;align-items:center;gap:6px">
              <el-icon :size="16"><component :is="getNodeIcon({icon:ico,type:'menu'})" /></el-icon>
              <span>{{ ico || '无' }}</span>
            </span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="父菜单"><el-tree-select v-model="form.parent_id" :data="treeData" :props="{ children:'children', label:'name', value:'id' }" placeholder="无（根节点）" clearable check-strictly filterable style="width:100%" /></el-form-item>
    </el-form>
    <template #footer><el-button @click="createVisible = false">取消</el-button><el-button type="primary" @click="handleCreate">确定</el-button></template>
  </el-dialog>

  <el-dialog v-model="editVisible" title="编辑菜单" width="420px" :close-on-click-modal="false">
    <el-form :model="editForm" label-width="80px">
      <el-form-item label="标识"><el-input v-model="editForm.name" /></el-form-item>
      <el-form-item label="描述"><el-input v-model="editForm.description" /></el-form-item>
      <el-form-item label="类型"><el-select v-model="editForm.type" style="width:100%"><el-option v-for="t in typeOptions" :key="t.value" :label="t.label" :value="t.value" /></el-select></el-form-item>
      <el-form-item label="图标">
        <el-select v-model="editForm.icon" filterable clearable placeholder="选择图标（仅目录/菜单）" style="width:100%">
          <el-option v-for="ico in iconOptions" :key="ico" :label="ico || '无'" :value="ico">
            <span style="display:flex;align-items:center;gap:6px">
              <el-icon :size="16"><component :is="getNodeIcon({icon:ico,type:'menu'})" /></el-icon>
              <span>{{ ico || '无' }}</span>
            </span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="父菜单"><el-tree-select v-model="editForm.parent_id" :data="treeData" :props="{ children:'children', label:'name', value:'id' }" placeholder="无（根节点）" clearable check-strictly filterable style="width:100%" /></el-form-item>
    </el-form>
    <template #footer><el-button @click="editVisible = false">取消</el-button><el-button type="primary" @click="handleEdit">保存</el-button></template>
  </el-dialog>
</template>

<style scoped>
/* ====== 页面卡片 ====== */
.page-card {
  border-radius: 8px;
}
.page-card :deep(.el-card__body) {
  padding: 16px 20px;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}
.header-actions {
  display: flex;
  gap: 8px;
}

/* ====== 工具栏 ====== */
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 6px;
}
.toolbar-btn {
  padding: 5px 8px;
}
.node-count {
  color: #909399;
  font-size: 13px;
  margin-left: 6px;
  white-space: nowrap;
}

/* ====== 类型图例 ====== */
.type-legend {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 10px 0 12px;
  color: #606266;
  font-size: 13px;
  border-bottom: 1px solid #ebeef5;
  margin-bottom: 4px;
}
.legend-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  cursor: default;
}
.legend-hint {
  margin-left: auto;
  color: #c0c4cc;
  font-size: 12px;
  font-style: italic;
}

/* ====== 树节点 ====== */
.menu-tree {
  font-size: 14px;
  margin-top: 4px;
}
:deep(.el-tree-node__content) {
  height: 38px;
  border-radius: 6px;
  transition: all 0.15s ease;
  padding: 0 4px;
}
:deep(.el-tree-node__content:hover) {
  background-color: #f5f7fa;
}
:deep(.el-tree-node.is-current > .el-tree-node__content) {
  background-color: #ecf5ff;
  color: #409eff;
}
:deep(.el-tree-node.is-current .node-name) {
  color: #409eff;
}

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 2px 0;
}
.node-info {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}
.node-icon {
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
}
.node-name {
  font-weight: 500;
  white-space: nowrap;
  transition: color 0.15s;
}
.node-desc {
  color: #999;
  margin-left: 2px;
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 200px;
}
.node-icon-name {
  color: #aaa;
  font-size: 11px;
  margin-left: 4px;
  font-style: italic;
}
.node-tag {
  flex-shrink: 0;
  margin-left: 8px;
}

/* ====== 右键菜单 ====== */
.ctx-menu {
  position: fixed;
  z-index: 9999;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  padding: 4px 0;
  min-width: 140px;
}
.ctx-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  cursor: pointer;
  font-size: 13px;
  color: #333;
  user-select: none;
  transition: background 0.1s;
}
.ctx-item:hover {
  background: #ecf5ff;
  color: #409eff;
}
.ctx-danger:hover {
  background: #fef0f0;
  color: #f56c6c;
}
.ctx-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  font-size: 14px;
}
.ctx-divider {
  height: 1px;
  background: #ebeef5;
  margin: 4px 0;
}

/* 右键菜单动画 */
.ctx-fade-enter-active {
  transition: opacity 0.15s ease, transform 0.12s ease;
}
.ctx-fade-leave-active {
  transition: opacity 0.1s ease;
}
.ctx-fade-enter-from {
  opacity: 0;
  transform: scale(0.95);
}
.ctx-fade-leave-to {
  opacity: 0;
}

/* 树展开折叠箭头美化 */
:deep(.el-tree-node__expand-icon) {
  font-size: 14px;
  color: #c0c4cc;
  transition: transform 0.2s;
}
:deep(.el-tree-node__expand-icon.is-leaf) {
  color: transparent;
}
</style>
