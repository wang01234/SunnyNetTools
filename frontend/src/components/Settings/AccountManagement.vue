<template>
  <div class="account-management">
    <div class="toolbar">
      <el-button type="primary" @click="openAddDialog">添加账号</el-button>
    </div>
    
    <el-table :data="accounts" style="width: 100%; margin-top: 20px" border>
      <el-table-column prop="username" label="用户名" min-width="200">
        <template #default="scope">
          <div style="display: flex; align-items: center;">
            <span>{{ scope.row.username }}</span>
            <el-tag v-if="scope.row.isDefault" size="small" type="info" style="margin-left: 8px">默认</el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" align="center">
        <template #default="scope">
          <el-button link type="primary" size="small" @click="openEditDialog(scope.row)">修改密码</el-button>
          <el-button 
            link 
            type="danger" 
            size="small" 
            :disabled="scope.row.isDefault"
            @click="deleteAccount(scope.row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加账号对话框 -->
    <el-dialog v-model="addDialogVisible" title="添加账号" width="400px">
      <el-form :model="addForm" :rules="addRules" ref="addFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="addForm.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="addForm.password" type="password" placeholder="请输入密码" show-password></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAdd" :loading="addLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 修改密码对话框 -->
    <el-dialog v-model="editDialogVisible" title="修改密码" width="400px">
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="editForm.username" disabled></el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input v-model="editForm.password" type="password" placeholder="请输入新密码" show-password></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEdit" :loading="editLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { CallGoDo } from '../CallbackEventsOn.js'

export default {
  name: 'AccountManagement',
  setup() {
    const accounts = ref([])
    const defaultUsername = 'admin'
    
    const addDialogVisible = ref(false)
    const editDialogVisible = ref(false)
    const addLoading = ref(false)
    const editLoading = ref(false)
    
    const addFormRef = ref(null)
    const editFormRef = ref(null)
    
    const addForm = ref({
      username: '',
      password: ''
    })
    
    const editForm = ref({
      username: '',
      password: ''
    })
    
    const addRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' }
      ]
    }
    
    const editRules = {
      password: [
        { required: true, message: '请输入新密码', trigger: 'blur' }
      ]
    }

    const loadAccounts = async () => {
      try {
        const currentUser = window.localStorage.getItem('currentUser') || ''
        const result = await CallGoDo('获取账号列表', { currentUser })
        if (result) {
          accounts.value = result.map(username => ({
            username,
            isDefault: username === defaultUsername
          }))
        }
      } catch (error) {
        ElMessage.error('获取账号列表失败')
      }
    }

    const openAddDialog = () => {
      addForm.value = { username: '', password: '' }
      addDialogVisible.value = true
    }

    const submitAdd = async () => {
      if (!addFormRef.value) return
      
      await addFormRef.value.validate(async (valid) => {
        if (!valid) return
        
        addLoading.value = true
        try {
          const result = await CallGoDo('添加账号', {
            username: addForm.value.username,
            password: addForm.value.password
          })
          
          if (result.success === true) {
            ElMessage.success('添加成功')
            addDialogVisible.value = false
            loadAccounts()
          } else {
            ElMessage.error(result.message || '添加失败')
          }
        } catch (error) {
          ElMessage.error('添加失败')
        } finally {
          addLoading.value = false
        }
      })
    }

    const openEditDialog = (row) => {
      editForm.value = {
        username: row.username,
        password: ''
      }
      editDialogVisible.value = true
    }

    const submitEdit = async () => {
      if (!editFormRef.value) return
      
      await editFormRef.value.validate(async (valid) => {
        if (!valid) return
        
        editLoading.value = true
        try {
          const result = await CallGoDo('修改密码', {
            username: editForm.value.username,
            password: editForm.value.password
          })
          
          if (result.success === true) {
            ElMessage.success('修改成功')
            editDialogVisible.value = false
          } else {
            ElMessage.error(result.message || '修改失败')
          }
        } catch (error) {
          ElMessage.error('修改失败')
        } finally {
          editLoading.value = false
        }
      })
    }

    const deleteAccount = (row) => {
      ElMessageBox.confirm(
        `确定要删除账号 "${row.username}" 吗？`,
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(async () => {
        try {
          const result = await CallGoDo('删除账号', {
            username: row.username
          })
          
          if (result.success === true) {
            ElMessage.success('删除成功')
            loadAccounts()
          } else {
            ElMessage.error(result.message || '删除失败')
          }
        } catch (error) {
          ElMessage.error('删除失败')
        }
      }).catch(() => {})
    }

    onMounted(() => {
      loadAccounts()
    })

    return {
      accounts,
      addDialogVisible,
      editDialogVisible,
      addLoading,
      editLoading,
      addFormRef,
      editFormRef,
      addForm,
      editForm,
      addRules,
      editRules,
      openAddDialog,
      submitAdd,
      openEditDialog,
      submitEdit,
      deleteAccount
    }
  }
}
</script>

<style scoped>
.account-management {
  padding: 10px;
}

.toolbar {
  margin-bottom: 10px;
}
</style>
