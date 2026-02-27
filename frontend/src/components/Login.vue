<script>
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { CallGoDo } from './CallbackEventsOn.js'
import { WindowMinimise, WindowMaximise, WindowUnmaximise, Quit } from '../../wailsjs/runtime/runtime'

export default {
  name: 'Login',
  data() {
    return {
      loginForm: reactive({
        username: '',
        password: ''
      }),
      rememberMe: false,
      loginRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 1, message: '密码不能为空', trigger: 'blur' }
        ]
      },
      loading: false,
      isMaximised: false
    }
  },
  mounted() {
    // 检查是否保存了登录信息
    this.loadRememberedCredentials()
  },
  methods: {
    // 加载记住的凭证
    loadRememberedCredentials() {
      const savedUsername = localStorage.getItem('rememberedUsername')
      const savedPassword = localStorage.getItem('rememberedPassword')
      const savedRememberMe = localStorage.getItem('rememberMe')
      
      if (savedRememberMe === 'true' && savedUsername && savedPassword) {
        this.loginForm.username = savedUsername
        this.loginForm.password = savedPassword
        this.rememberMe = true
        // 自动登录
        this.handleLogin()
      }
    },
    // 保存记住的凭证
    saveRememberedCredentials() {
      if (this.rememberMe) {
        localStorage.setItem('rememberedUsername', this.loginForm.username)
        localStorage.setItem('rememberedPassword', this.loginForm.password)
        localStorage.setItem('rememberMe', 'true')
      } else {
        localStorage.removeItem('rememberedUsername')
        localStorage.removeItem('rememberedPassword')
        localStorage.setItem('rememberMe', 'false')
      }
    },
    async handleLogin() {
      this.$refs.loginForm.validate(async (valid) => {
        if (!valid) {
          return false
        }
        
        this.loading = true
        try {
          // 调用后端登录验证 API
          const result = await CallGoDo('用户登录', {
            username: this.loginForm.username,
            password: this.loginForm.password
          })
          
          if (result.success === true) {
            // 保存记住的凭证
            this.saveRememberedCredentials()
            ElMessage({
              message: '登录成功',
              type: 'success'
            })
            // 存储当前用户名和角色
            window.localStorage.setItem('currentUser', result.username || this.loginForm.username)
            window.localStorage.setItem('currentRole', result.role || 'user')
            // 触发登录成功事件
            window.dispatchEvent(new CustomEvent('login-success'))
          } else {
            ElMessage({
              message: result.message || '用户名或密码错误',
              type: 'error'
            })
          }
        } catch (error) {
          ElMessage({
            message: '登录失败，请稍后重试',
            type: 'error'
          })
        } finally {
          this.loading = false
        }
      })
    },
    resetForm() {
      this.$refs.loginForm.resetFields()
    },
    handleMinimize() {
      WindowMinimise()
    },
    handleMaximize() {
      if (this.isMaximised) {
        WindowUnmaximise()
        this.isMaximised = false
      } else {
        WindowMaximise()
        this.isMaximised = true
      }
    },
    handleClose() {
      Quit()
    }
  }
}
</script>

<template>
  <div class="login-container">
    <div class="window-controls">
      <button class="control-btn minimize" @click="handleMinimize" title="最小化">
        <span>&#x2212;</span>
      </button>
      <button class="control-btn maximize" @click="handleMaximize" title="最大化">
        <span>{{ isMaximised ? '&#x2750;' : '&#x25A1;' }}</span>
      </button>
      <button class="control-btn close" @click="handleClose" title="关闭">
        <span>&#x2715;</span>
      </button>
    </div>
    <div class="login-box">
      <div class="login-header">
        <h2>大神工具 登录</h2>
      </div>
      <el-form
        ref="loginForm"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        autocomplete="on"
        label-position="left"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名"
            name="username"
            type="text"
            tabindex="1"
            autocomplete="on"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            placeholder="密码"
            name="password"
            type="password"
            tabindex="2"
            autocomplete="on"
            prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item>
          <el-checkbox v-model="rememberMe">记住我</el-checkbox>
        </el-form-item>
        
        <el-button
          :loading="loading"
          type="primary"
          size="large"
          style="width: 100%; margin-bottom: 30px"
          @click.prevent="handleLogin"
        >
          {{ loading ? '登录中...' : '登 录' }}
        </el-button>
        
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100%;
  width: 100%;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  position: relative;
}

.window-controls {
  position: absolute;
  top: 0;
  right: 0;
  display: flex;
  -webkit-app-region: no-drag;
}

.control-btn {
  width: 46px;
  height: 32px;
  border: none;
  background: transparent;
  color: #fff;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.control-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.control-btn.close:hover {
  background: #ff5f57;
}

.control-btn span {
  line-height: 1;
}

.login-box {
  width: 400px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 10px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
  margin: 0;
  color: #333;
  font-size: 24px;
  font-weight: 600;
}

.login-form {
  margin-top: 20px;
}

.tips {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
}

:deep(.el-input__wrapper) {
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #409eff inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #409eff inset;
}

:deep(.el-form-item) {
  margin-bottom: 24px;
}

:deep(.el-button--primary) {
  background-color: #409eff;
  border-color: #409eff;
}

:deep(.el-button--primary:hover) {
  background-color: #66b1ff;
  border-color: #66b1ff;
}
</style>
