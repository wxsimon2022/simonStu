<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../utils/request'
import { ElMessage } from 'element-plus'

const profile = ref({ username: '', real_name: '', phone: '', email: '', create_time: '' })
const editForm = ref({ real_name: '', phone: '', email: '' })
const password = ref('')
const confirm = ref('')
const loading = ref(false)
const saving = ref(false)

async function fetchProfile() {
  loading.value = true
  const res = await api('/api/auth/profile')
  if (res.code === 200) {
    const d = res.data
    profile.value = d
    editForm.value = { real_name: d.real_name || '', phone: d.phone || '', email: d.email || '' }
  } else {
    ElMessage.error(res.message || '加载个人信息失败')
  }
  loading.value = false
}

async function handleSave() {
  // 密码校验
  if (password.value || confirm.value) {
    if (password.value.length < 6) { ElMessage.warning('密码长度不少于 6 位'); return }
    if (password.value !== confirm.value) { ElMessage.warning('两次输入的密码不一致'); return }
  }

  const body = {}
  if (editForm.value.real_name !== profile.value.real_name) body.real_name = editForm.value.real_name
  if (editForm.value.phone !== profile.value.phone) body.phone = editForm.value.phone
  if (editForm.value.email !== profile.value.email) body.email = editForm.value.email
  if (password.value) body.password = password.value

  if (!Object.keys(body).length) {
    ElMessage.info('没有需要保存的更改')
    return
  }

  saving.value = true
  const res = await api('/api/auth/profile', { method: 'PUT', body: JSON.stringify(body) })
  saving.value = false

  if (res.code === 200) {
    ElMessage.success('保存成功')
    profile.value = { ...profile.value, ...body }
    password.value = ''
    confirm.value = ''
  } else {
    ElMessage.error(res.message || '保存失败')
  }
}

onMounted(fetchProfile)
</script>

<template>
  <div class="profile-page">
    <el-card shadow="never" class="section-card">
      <template #header>
        <span class="section-title">个人中心</span>
      </template>
      <el-form :model="editForm" label-width="80px" v-loading="loading">
        <!-- 只读信息 -->
        <el-form-item label="用户名">
          <el-input :model-value="profile.username" disabled class="readonly-input" />
        </el-form-item>
        <el-form-item label="创建时间">
          <el-input :model-value="profile.create_time" disabled class="readonly-input" />
        </el-form-item>

        <el-divider />

        <!-- 可编辑信息 -->
        <el-form-item label="姓名">
          <el-input v-model="editForm.real_name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="editForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="editForm.email" placeholder="请输入邮箱" />
        </el-form-item>

        <el-divider />

        <!-- 密码 -->
        <el-form-item label="新密码">
          <el-input v-model="password" type="password" placeholder="留空不修改，不少于 6 位" show-password />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="confirm" type="password" placeholder="再次输入新密码" show-password />
        </el-form-item>

        <!-- 提交按钮 -->
        <el-form-item>
          <el-button type="primary" size="large" :loading="saving" @click="handleSave" class="save-btn">
            保存修改
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.profile-page {
  max-width: 560px;
  margin: 0 auto;
  padding: 0;
}
.section-card {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
.section-card :deep(.el-card__body) {
  padding: 20px 28px;
}
.section-card :deep(.el-card__header) {
  padding: 14px 28px;
  border-bottom: 1px solid #e5e6eb;
}
.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
}
.readonly-input :deep(.el-input__wrapper) {
  background: #f7f8fa;
  box-shadow: 0 0 0 1px #e4e7ed inset;
}
.readonly-input :deep(.el-input__inner) {
  color: #86909c;
}
:deep(.el-divider) {
  margin: 16px 0;
}
:deep(.el-form-item) {
  margin-bottom: 18px;
}
:deep(.el-form-item:last-child) {
  margin-bottom: 0;
}
.save-btn {
  width: 100%;
}
</style>
