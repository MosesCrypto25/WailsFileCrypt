<template>
    <div class="container">
      <button @click="pickFile">选择文件</button>
      <input v-model="key" placeholder="输入16/24/32字节密钥" />
      <button @click="encrypt">加密</button>
      <button @click="decrypt">解密</button>
      <p v-if="result">{{ result }}</p>
      <p v-if="error" class="error">{{ error }}</p>
      <!-- 进度条显示 -->
    <div class="progress-container" v-show="showProgress">
      <div 
        class="progress-bar" 
        :style="{ width: progress + '%' }"
      ></div>
      <span class="status-text">{{ statusText }}</span>
    </div>
      <!-- 隐藏的本地文件选择器 -->
      <input 
        type="file"
        ref="fileInput"
        @change="handleFileSelect"
        style="display: none;"
      />
    </div>
  </template>
   
  <script>
  export default {
    data() {
      return {
      selectedFilePath: '',
      key: '',
      result: '',
      error: '',
      progress: 0,
      statusText: '',
      showProgress: false,
    }
    },

  mounted() {
    // 监听后端进度事件 
    window.runtime.EventsOn('progress',  (data) => {
      this.progress  = data.percentage  
      this.statusText  = this.getStatusText(data.status) 
      this.showProgress  = true 
      
      if (data.status  === 'done') {
        setTimeout(() => {
          this.showProgress  = false 
          this.progress  = 0 
        }, 2000)
      }
    })
  },
    methods: {

      async pickFile() {
      try {
        const filePath = await window.go.main.App.GetFilePath(); 
        this.selectedFilePath  = filePath; // 直接获取后端返回的路径 
      } catch (err) {
        console.error(" 文件选择失败:", err);
      }
    },

      // 打开系统级文件对话框 
      async openFileDialog() {
        try {
          const path = await window.go.main.App.OpenFileDialog(); 
          if (path) {
            this.selectedFilePath  = path;
            this.error  = '';
          }
        } catch (err) {
          this.error  = `文件选择失败: ${err}`;
        }
      },
   
      // 本地文件选择器（备选方案）
      handleFileSelect(e) {
        const file = e.target.files[0]; 
        this.selectedFilePath  = file ? file.path  : '';
        this.$refs.fileInput.value  = null; // 清空选择器 
      },
   
      async encrypt() {
        if (!this.validateKey())  return;
        try {
          const result = await window.go.main.App.EncryptFile( 
            this.selectedFilePath,  
            this.key  
          );
          this.result  = `加密文件已保存至: ${result}`;
        } catch (err) {
          this.error  = `加密失败: ${err}`;
        }
      },
    async decrypt() {
      if (!this.validateKey())  return 
      
      try {
        const result = await window.go.main.App.DecryptFile( 
          this.selectedFilePath, 
          this.key  
        )
        this.result  = `解密成功：${result}`
        this.error  = ''
      } catch (err) {
        this.error  = `解密失败：${err.message}` 
      }
    },
      validateKey() {
        if (this.key.length >= 6) {
            return true;
        }
        this.error = '密钥长度必须至少为6字节';
        return false;
      },

      getStatusText(status) {
      const statusMap = {
        encrypting: '正在加密...',
        decrypting: '正在解密...',
        done: '处理完成',
        error: '处理出错',
      }
      return statusMap[status] || ''
    }
    }
  }
  </script>
  <style scoped>
  .container {
    padding: 20px;
    max-width: 600px;
    margin: 0 auto;
  }
   
  button {
    padding: 10px 20px;
    margin-right: 10px;
    cursor: pointer;
  }
   
  input {
    padding: 8px;
    width: 200px;
    margin-bottom: 10px;
  }
   
  .error {
    color: red;
  }
   
  .progress-container {
    margin-top: 20px;
    padding: 15px;
    background-color: #f5f5f5;
    border-radius: 5px;
  }
   
  .progress-bar {
    height: 20px;
    background-color: #4CAF50;
    border-radius: 3px;
    transition: width 0.3s ease;
  }
   
  .status-text {
    display: block;
    margin-top: 5px;
    color: #333;
  }
  </style>