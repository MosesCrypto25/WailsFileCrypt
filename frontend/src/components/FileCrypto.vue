<template>
  <div class="container">
    <button @click="pickFile">选择文件</button>
    <input v-model="key" placeholder="输入6位以上密码" />
    <button @click="encrypt">加密</button>
    <button @click="decrypt">解密</button>
    <p v-if="selectedFilePath" class="info-box file-path">文件路径: {{ selectedFilePath }}</p>
    <p v-if="result" class="info-box result">加密/解密结果: {{ result }}</p>
    <p v-if="error" class="error">{{ error }}</p>
    <!-- 进度条显示 -->
    <div class="progress-container" v-show="showProgress">
      <div
        class="progress-bar"
        :style="{ width: progress + '%' }"
      ></div>
      <span class="status-text">{{ statusText }}</span>
    </div>
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
 
.info-box {
  padding: 10px;
  border-radius: 5px;
  margin-top: 10px;
  color: #000;
}

.file-path {
  background-color: #e0f7fa; /* 淡蓝色底色 */
}

.result {
  background-color: #e8f5e9; /* 淡绿色底色 */
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