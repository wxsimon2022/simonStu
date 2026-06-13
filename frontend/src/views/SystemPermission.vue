<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { api } from '../utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { Folder, FolderOpened, Menu, SwitchButton, Search, Edit, Plus, Delete, Refresh, Close } from '@element-plus/icons-vue'

const treeRef = ref(null)
const treeData = ref([])
const loading = ref(false)
const selectedNode = ref(null)
const filterText = ref('')
const expandedKeys = ref([])
const createVisible = ref(false)
const editVisible = ref(false)
const contextMenu = ref({ visible: false, x: 0, y: 0 })
const form = ref({ name: '', description: '', type: 'menu', icon: '', parent_id: null })
const editForm = ref({ id: 0, name: '', description: '', type: 'menu', icon: '', parent_id: null })

const treeProps = { children: 'children', label: 'name' }
// 常用 Element Plus 图标列表
const iconOptions = [
  '', 'User', 'UserFilled', 'Management', 'Setting', 'Menu',
  'Folder', 'FolderOpened', 'HomeFilled', 'Tools', 'Operation',
  'Histogram', 'DataAnalysis', 'Aim', 'Bell',
  'ChatDotSquare', 'CircleCheck', 'Document', 'Notebook',
  'Files', 'Search', 'Plus', 'Edit', 'Delete', 'Link',
  'SwitchButton', 'Refresh', 'Share', 'Upload',
  'Flag', 'Star', 'Collection', 'List',
]

const typeOptions = [
  { label: '目录', value: 'dir' },
  { label: '菜单', value: 'menu' },
  { label: '按钮', value: 'btn' },
]

// 获取节点深度
function getNodeDepth(nodes, targetId, depth = 0) {
  for (const n of nodes) {
    if (n.id === targetId) return depth
    if (n.children?.length) {
      const d = getNodeDepth(n.children, targetId, depth + 1)
      if (d !== -1) return d
    }
  }
  return -1
}

function typeLabel(t) { return { dir: '目录', menu: '菜单', btn: '按钮' }[t] || '菜单' }
function typeTagType(t) { return { dir: 'warning', menu: 'primary', btn: 'info' }[t] || 'primary' }

// 树节点图标：优先用数据库中配置的 icon，否则按 type 显示默认图标
function getNodeIcon(data) {
  if (data.icon && ElementPlusIconsVue[data.icon]) return ElementPlusIconsVue[data.icon]
  const fallback = { dir: Folder, menu: Menu, btn: SwitchButton }
  return fallback[data.type] || Menu
}

// 计算树节点总数
function countNodes(nodes) {
  let c = 0
  for (const n of nodes) { c++; if (n.children?.length) c += countNodes(n.children) }
  return c
}
const flatCount = computed(() => countNodes(treeData.value))
const selectedDepth = computed(() => {
  if (!selectedNode.value) return -1
  return getNodeDepth(treeData.value, selectedNode.value.id)
})

// 搜索高亮 — 将匹配文本包裹在 <span class="highlight"> 中
const highlightMatch = (text, query) => {
  if (!query || !text) return text
  const q = query.toLowerCase()
  const t = String(text)
  const idx = t.toLowerCase().indexOf(q)
  if (idx === -1) return t
  return t.slice(0, idx) +
    '<span class="highlight">' + t.slice(idx, idx + q.length) + '</span>' +
    t.slice(idx + q.length)
}

watch(filterText, (v) => { treeRef.value?.filter(v) })
function filterNode(value, data, node) {
  if (!value) return true
  const q = value.toLowerCase()
  const matched = data.name.toLowerCase().includes(q) || (data.description || '').toLowerCase().includes(q)
  // 如果当前节点匹配，展开其所有祖先节点
  if (matched && treeRef.value) {
    let parent = node?.parent
    while (parent && parent.key !== undefined) {
      expandedKeys.value.push(parent.key)
      parent = parent.parent
    }
  }
  return matched
}

function setExpand(expand) {
  if (expand) {
    const keys = []
    function walk(nodes) {
      for (const n of nodes) {
        keys.push(n.id)
        if (n.children?.length) walk(n.children)
      }
    }
    walk(treeData.value)
    expandedKeys.value = keys
  } else {
    expandedKeys.value = []
  }
}


function onKeyDown(e) {
  if (e.key === 'Enter' && selectedNode.value && !createVisible.value && !editVisible.value) openEdit()
  if ((e.key === 'Delete' || e.key === 'Backspace') && selectedNode.value && !createVisible.value && !editVisible.value) handleDelete()
}
function closeCtx() { contextMenu.value.visible = false }
function onBodyClick() {
  closeCtx()
  // 点击空白区域取消选中（但保留点击树节点本身触发的选中）
}

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

// 拖拽排序
const dragging = ref(false)

function buildReorderPayload(nodes, parentId) {
  let items = []
  for (let i = 0; i < nodes.length; i++) {
    const n = nodes[i]
    items.push({ id: n.id, parent_id: parentId ?? null, sort_order: i })
    if (n.children?.length) {
      items = items.concat(buildReorderPayload(n.children, n.id))
    }
  }
  return items
}

async function handleDrop(draggingNode, dropNode, dropType, ev) {
  dragging.value = false
  const items = buildReorderPayload(treeData.value)
  if (!items.length) return
  const res = await api('/api/permission/reorder', {
    method: 'PUT',
    body: JSON.stringify({ items }),
  })
  if (res.code === 200) {
    ElMessage.success('排序已更新')
  } else {
    ElMessage.error(res.message || '排序更新失败')
    fetchData()
  }
}

function allowDrop(draggingNode, dropNode, type) {

function onDragStart() { dragging.value = true }
function onDragEnd() { dragging.value = false }

  return ['prev', 'next', 'inner'].includes(type)
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

function onEditNode(data) {
  selectedNode.value = data
  contextMenu.value.visible = false
  editForm.value = {
    id: data.id, name: data.name,
    description: data.description || '', type: data.type || 'menu',
    icon: data.icon || '',
    parent_id: data.parent_id ?? null,
  }
  editVisible.value = true
}

function onAddChild(data) {
  selectedNode.value = data
  contextMenu.value.visible = false
  form.value = { name: '', description: '', type: 'menu', icon: '', parent_id: data.id }
  createVisible.value = true
}

// 点击行内删除按钮时调用
function onDeleteNode(data) {
  selectedNode.value = data
  contextMenu.value.visible = false
  handleDelete()
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

    <!-- 选中节点信息栏 -->
    <transition name="info-slide">
      <div v-if="selectedNode" class="detail-bar">
        <div class="detail-left">
          <el-icon :size="16"><component :is="getNodeIcon(selectedNode)" /></el-icon>
          <span class="detail-name">{{ selectedNode.name }}</span>
          <span v-if="selectedNode.description" class="detail-desc">—— {{ selectedNode.description }}</span>
        </div>
        <div class="detail-right">
          <el-tag :type="typeTagType(selectedNode.type)" size="small" effect="plain">{{ typeLabel(selectedNode.type) }}</el-tag>
          <span class="detail-depth">L{{ selectedDepth + 1 }}</span>
          <span class="detail-id">#{{ selectedNode.id }}</span>
          <el-tooltip content="取消选中">
            <el-button size="small" text @click.stop="selectedNode = null"><el-icon><Close /></el-icon></el-button>
          </el-tooltip>
        </div>
      </div>
    </transition>

    <div class="toolbar">
      <el-input v-model="filterText" placeholder="搜索菜单名称或描述…" clearable :prefix-icon="Search" size="small" style="width:260px" />
      <el-tooltip content="刷新菜单树">
        <el-button size="small" :loading="loading" class="toolbar-btn" @click="fetchData">
          <el-icon><Refresh /></el-icon>
        </el-button>
      </el-tooltip>
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
      <span class="legend-hint">拖拽排序 · 回车编辑 · Delete 删除</span>
      <span class="legend-hint-mobile">拖拽排序 · 点击选中 · 右键操作</span>
    </div>

    <el-tree
      ref="treeRef"
      :data="treeData"
      :props="treeProps"
      :filter-node-method="filterNode"
      node-key="id"
      highlight-current
      v-model:expanded-keys="expandedKeys"
      draggable
      :allow-drop="allowDrop"
      @node-drop="handleDrop"
      @node-drag-start="() => dragging = true"
      @node-drag-end="() => dragging = false"
      v-loading="loading"
      @node-click="onNodeClick"
      @node-contextmenu="onContextMenu"
      class="menu-tree"
    >
      <template #default="{ data }">
        <div class="tree-node">
          <div class="node-content">
            <div class="node-info">
              <el-icon :size="16" class="node-icon">
                <component :is="getNodeIcon(data)" />
              </el-icon>
              <span class="node-name" v-html="filterText ? highlightMatch(data.name, filterText) : data.name"></span>
              <span v-if="data.description" class="node-desc" v-html="filterText ? highlightMatch(data.description, filterText) : data.description"></span>
              <span v-if="data.icon" class="node-icon-name">{{ data.icon }}</span>
            </div>
            <el-tag :type="typeTagType(data.type)" size="small" effect="plain" class="node-tag">
              {{ typeLabel(data.type) }}
            </el-tag>
          </div>
          <div class="node-actions">
            <el-tooltip content="编辑">
              <el-button size="small" link type="primary" @click.stop="onEditNode(data)">
                <el-icon><Edit /></el-icon>
              </el-button>
            </el-tooltip>
            <el-tooltip content="新增子菜单">
              <el-button size="small" link type="primary" @click.stop="onAddChild(data)">
                <el-icon><Plus /></el-icon>
              </el-button>
            </el-tooltip>
            <el-tooltip content="删除">
              <el-button size="small" link type="danger" @click.stop="onDeleteNode(data)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </el-tooltip>
          </div>
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
          <span class="ctx-icon">&check;</span> 编辑
        </div>
        <div class="ctx-divider"></div>
        <div class="ctx-item ctx-danger" @click="handleDelete">
          <span class="ctx-icon">&times;</span> 删除
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
        <el-select v-model="form.icon" filterable clearable placeholder="选择图标（仅目录/菜单）">
          <template #prefix><el-icon v-if="form.icon && ElementPlusIconsVue[form.icon]" :size="16"><component :is="getNodeIcon({icon:form.icon,type:'menu'})" /></el-icon></template>
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
      <div class="dialog-id-hint">ID: {{ editForm.id }}</div>
      <el-form-item label="标识"><el-input v-model="editForm.name" /></el-form-item>
      <el-form-item label="描述"><el-input v-model="editForm.description" /></el-form-item>
      <el-form-item label="类型"><el-select v-model="editForm.type" style="width:100%"><el-option v-for="t in typeOptions" :key="t.value" :label="t.label" :value="t.value" /></el-select></el-form-item>
      <el-form-item label="图标">
        <el-select v-model="editForm.icon" filterable clearable placeholder="选择图标（仅目录/菜单）">
          <template #prefix><el-icon v-if="editForm.icon && ElementPlusIconsVue[editForm.icon]" :size="16"><component :is="getNodeIcon({icon:editForm.icon,type:'menu'})" /></el-icon></template>
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
  box-shadow: 0 1px 3px rgba(0, 0, 0, .04);
}
.page-card :deep(.el-card__body) {
  padding: 16px 24px;
}
.page-card :deep(.el-card__header) {
  padding: 14px 24px;
  border-bottom: 1px solid #e5e6eb;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
}
.header-actions {
  display: flex;
  gap: 8px;
}

/* ====== 选中节点详情栏 ====== */
.detail-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  margin-bottom: 12px;
  background: linear-gradient(135deg, #e8f4ff 0%, #f0f8ff 100%);
  border: 1px solid #b3d8ff;
  border-radius: 6px;
  font-size: 13px;
  min-height: 32px;
}
.detail-left {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
  flex: 1;
}
.detail-name {
  font-weight: 600;
  color: #1d2129;
}
.detail-desc {
  color: #86909c;
  font-size: 12px;
  margin-left: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.detail-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
  margin-left: 12px;
}
.detail-depth {
  font-size: 11px;
  color: #86909c;
  background: #f2f3f5;
  padding: 0 6px;
  border-radius: 3px;
  line-height: 18px;
  font-variant-numeric: tabular-nums;
  font-weight: 500;
}
.detail-id {
  font-size: 11px;
  color: #c9cdd4;
  font-family: monospace;
}

/* ====== 工具栏 ====== */
.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}
.toolbar :deep(.el-input__wrapper) {
  border-radius: 6px;
  box-shadow: 0 0 0 1px #e4e7ed inset;
  transition: box-shadow 0.2s;
}
.toolbar :deep(.el-input__wrapper:hover),
.toolbar :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #409eff inset;
}
.toolbar :deep(.el-input__prefix-inner) .el-icon {
  color: #c9cdd4;
}
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: auto;
}
.toolbar-btn {
  padding: 5px 8px;
  border-radius: 6px;
  transition: all 0.15s;
}
.toolbar-btn:hover {
  background: #e8f4ff;
  color: #409eff;
}
.node-count {
  color: #86909c;
  font-size: 12px;
  margin-left: 8px;
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
}

/* ====== 类型图例 ====== */
.type-legend {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 0 10px;
  color: #4e5969;
  font-size: 13px;
  border-bottom: 1px solid #e5e6eb;
  margin-bottom: 0;
}
.legend-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  cursor: default;
  padding: 2px 6px;
  border-radius: 4px;
  background: #f7f8fa;
  line-height: 22px;
}
.legend-item .el-icon {
  display: inline-flex;
}
.legend-hint {
  margin-left: auto;
  color: #c9cdd4;
  font-size: 12px;
}
.legend-hint-mobile {
  display: none;
}

/* ====== 树节点 ====== */
.menu-tree {
  font-size: 14px;
  margin-top: 4px;
}
.menu-tree :deep(.el-tree-node__content) {
  height: 40px;
  border-radius: 0;
  transition: all 0.15s ease;
  padding: 0 0 0 4px;
  margin: 0;
  position: relative;
  border-left: 3px solid transparent;
}
.menu-tree :deep(.el-tree-node__content:hover) {
  background-color: #f5f7fa;
}
.menu-tree :deep(.el-tree-node.is-current > .el-tree-node__content) {
  background: linear-gradient(90deg, #e8f4ff 0%, #f0f8ff 100%);
  border-left-color: #409eff;
  color: #409eff;
  font-weight: 500;
}
.menu-tree :deep(.is-current .node-name),
.menu-tree :deep(.is-current .node-info) {
  color: #409eff;
  font-weight: 500;
}
.menu-tree :deep(.is-current .node-icon) {
  color: #409eff;
}

/* ====== 层级缩进连线 ====== */
.menu-tree :deep(.el-tree-node__children) {
  position: relative;
}
.menu-tree :deep(.el-tree-node__children > .el-tree-node::before) {
  content: '';
  position: absolute;
  left: 16px;
  top: 0;
  bottom: 50%;
  width: 1px;
  background: #dcdfe6;
}
.menu-tree :deep(.el-tree-node__children > .el-tree-node:last-child)::before {
  bottom: auto;
  height: 40px;
}
.menu-tree :deep(.el-tree-node__children > .el-tree-node)::after {
  content: '';
  position: absolute;
  left: 16px;
  top: 20px;
  width: 12px;
  height: 1px;
  background: #dcdfe6;
}
.menu-tree :deep(.el-tree-node__children > .el-tree-node:last-child)::after {
  top: 0;
  height: 20px;
  width: 12px;
}

.tree-node {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 0;
  min-width: 0;
  gap: 8px;
  position: relative;
}
.node-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}
.node-info {
  display: flex;
  align-items: center;
  gap: 5px;
  min-width: 0;
  flex: 1;
}
.node-icon {
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  width: 18px;
  justify-content: center;
  color: #86909c;
  transition: color 0.15s;
}
.node-name {
  font-weight: 500;
  color: #1d2129;
  white-space: nowrap;
  transition: color 0.15s;
  font-size: 13.5px;
}
.node-desc {
  color: #86909c;
  font-size: 11.5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 140px;
  margin-left: 2px;
}
.node-icon-name {
  color: #c9cdd4;
  font-size: 10px;
  margin-left: 2px;
  font-style: italic;
  font-family: monospace;
}

/* ====== 行内按钮 ====== */
.node-actions {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 1px;
  margin-left: 4px;
  opacity: 0;
  transition: opacity 0.15s ease;
}
.tree-node:hover .node-actions,
.menu-tree :deep(.is-current) .node-actions {
  opacity: 1;
}
.node-actions :deep(.el-button) {
  width: 26px;
  height: 26px;
  border-radius: 6px;
  transition: all 0.15s;
}
.node-actions :deep(.el-button:hover) {
  background: #d9ecff;
  transform: scale(1.1);
}
.node-actions :deep(.el-button--danger:hover) {
  background: #ffece8;
  transform: scale(1.1);
}

/* ====== 类型标签 ====== */
.node-tag {
  flex-shrink: 0;
  border-radius: 4px;
  font-weight: 400;
  letter-spacing: 0.2px;
  font-size: 12px;
  padding: 0 6px;
}

/* ====== 搜索高亮 ====== */
:deep(.highlight) {
  background: #ffd666;
  color: #1d2129;
  padding: 0 2px;
  border-radius: 2px;
  font-weight: 600;
}

/* ====== 右键菜单 ====== */
.ctx-menu {
  position: fixed;
  z-index: 9999;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 6px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  padding: 3px 0;
  min-width: 140px;
  overflow: hidden;
}
.ctx-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  cursor: pointer;
  font-size: 12.5px;
  color: #1d2129;
  user-select: none;
  transition: background 0.1s;
}
.ctx-item:hover {
  background: #e8f4ff;
  color: #409eff;
}
.ctx-danger:hover {
  background: #fff1f0;
  color: #f53f3f;
}
.ctx-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  font-size: 12px;
  color: #86909c;
}
.ctx-divider {
  height: 1px;
  background: #f0f0f0;
  margin: 4px 0;
}

/* 右键菜单动画 */
.ctx-fade-enter-active {
  transition: opacity 0.1s ease, transform 0.08s ease;
}
.ctx-fade-leave-active {
  transition: opacity 0.06s ease;
}
.ctx-fade-enter-from {
  opacity: 0;
  transform: scale(0.96);
}
.ctx-fade-leave-to {
  opacity: 0;
}

/* 详情栏动画 */
.info-slide-enter-active {
  transition: all 0.2s ease;
}
.info-slide-leave-active {
  transition: all 0.15s ease;
}
.info-slide-enter-from,
.info-slide-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

/* 树展开折叠箭头美化 */
.menu-tree :deep(.el-tree-node__expand-icon) {
  font-size: 12px;
  color: #c9cdd4;
  transition: transform 0.2s ease;
}
.menu-tree :deep(.el-tree-node__expand-icon.is-leaf) {
  color: transparent;
}
.menu-tree :deep(.el-tree-node__expand-icon.expanded) {
  transform: rotate(90deg);
}

/* ====== 弹窗 ====== */
:deep(.el-dialog__body) {
  padding: 18px 24px;
}
:deep(.el-dialog .el-form-item) {
  margin-bottom: 16px;
}
:deep(.el-dialog .el-form-item:last-child) {
  margin-bottom: 0;
}
:deep(.el-dialog__header) {
  padding: 14px 24px;
  margin: 0;
  border-bottom: 1px solid #e5e6eb;
}
:deep(.el-dialog__footer) {
  padding: 12px 24px;
  border-top: 1px solid #e5e6eb;
}

/* 弹窗 ID 提示 */
.dialog-id-hint {
  font-size: 11px;
  color: #c9cdd4;
  margin-bottom: 12px;
  font-family: monospace;
  text-align: right;
}

/* ====== 空状态 ====== */
:deep(.el-empty) {
  padding: 40px 0;
}
:deep(.el-empty__description) {
  margin-top: 8px;
}
:deep(.el-empty__description p) {
  color: #86909c;
  font-size: 13px;
}

/* 图标选择器预览 */
:deep(.el-select .el-icon) {
  display: inline-flex;
  align-items: center;
}

/* ====== 响应式 ====== */

/* ====== 拖拽状态 ====== */
.menu-tree.dragging :deep(.el-tree-node__content) {
  cursor: grab;
}
.menu-tree.dragging :deep(.el-tree-node__content:active) {
  cursor: grabbing;
}
.menu-tree :deep(.el-tree-node.is-drop-inner) > .el-tree-node__content {
  background: #ecf5ff;
  border: 1px dashed #409eff;
  border-radius: 4px;
}
.menu-tree :deep(.el-tree-node.is-drop-not-allow) > .el-tree-node__content {
  opacity: 0.5;
}
@media (max-width: 768px) {
  .page-card :deep(.el-card__body) {
    padding: 12px 16px;
  }
  .card-title {
    font-size: 14px;
  }
  .header-actions :deep(.el-button) {
    padding: 5px 10px;
    font-size: 12px;
  }
  .toolbar {
    flex-wrap: wrap;
    gap: 8px;
  }
  .toolbar .el-input {
    width: 100% !important;
  }
  .toolbar-right {
    margin-left: 0;
  }
  .type-legend {
    flex-wrap: wrap;
    gap: 8px;
  }
  .legend-hint {
    display: none;
  }
  .legend-hint-mobile {
    display: inline;
    color: #c9cdd4;
    font-size: 12px;
    margin-left: auto;
  }
  .detail-bar {
    flex-wrap: wrap;
    gap: 6px;
  }
  .detail-right {
    margin-left: 0;
  }
  .node-desc {
    max-width: 80px;
  }
  :deep(.el-dialog) {
    width: 90vw !important;
  }
}
</style>
