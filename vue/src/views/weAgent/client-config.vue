<template>
    <div class="app-container">
        <div class="config-item">
            <span class="field-label">广告系统开关</span>
            <el-switch v-model="ad"
                       active-color="#13ce66"
                       inactive-color="#ff4949">
            </el-switch>
        </div>

        <div class="config-item">
            <span class="field-label">当前版本号(格式1.0.1)：</span>
            <el-input v-model="version"
                      style="width:150px;"
                      placeholder="请输入内容"></el-input>
        </div>

        <div class="setButton">
            <el-button type="primary"
                       @click="handleConfigSet">提交</el-button>
        </div>
    </div>
</template>

<script>
import waves from '@/directive/waves' // waves directive
import { ClientConfigGet, ClientConfigSet } from '@/api/weagent'

export default {
    name: 'ClientConfig',
    data () {
        return {
            ad: true,
            version: "",
        }
    },
    created () {
        this.getConfig()
    },
    methods: {
        // 获取配置
        getConfig () {
            ClientConfigGet().then(response => {
                this.ad = response.ad
                this.version = response.version
            })
        },
        handleConfigSet () {
            var params = {
                ad: this.ad,
                version: this.version
            }
            var data = JSON.stringify(params)

            ClientConfigSet(data).then(response => {
                this.$message({
                    message: '操作Success',
                    type: 'success'
                })
            })
        }
    }
}
</script>



<style scoped>
.config-item {
    margin-top: 15px;
}
.setButton {
    margin-top: 30px;
}
</style>
