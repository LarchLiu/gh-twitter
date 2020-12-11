<template>
  <div class="container">
    <div class="w1210">
      <div class="bd">
        <div class="aside-wrap">
          <aside-box
            title="Twitter"
            :need-fixed="true" 
            id-name="header"
          >
            <div v-if="usersList">
              <div
                v-for="user in usersList"
                :key="user"
              >
                <h3>{{ user }}</h3>
              </div>
            </div>
          </aside-box>
        </div>
        <twitter
          class="detail"
          :detail="usersData ? usersData[0] : {}"
        />
      </div>
    </div>
  </div>
</template>
<script>
import { ref,getCurrentInstance, onMounted, watch } from 'vue';
import userApi from '/@/api/twitter/user';
import AsideBox from '/@/components/AsideBox/index.vue';
import Twitter from '/@/components/Twitter/index.vue';


export default {
    components: { AsideBox, Twitter },
    setup () {

        const { ctx } = getCurrentInstance();
        const usersList = ref([])
        const usersData = ref([])

        const onExit = ()=>{
            ctx.$router.push({
                path: '/login'
            })
        }

        const getUserList = ()=> {
            userApi.getUsersData().then(res => {
                console.log(usersList.value)
                usersList.value = res.replace(/\s*/g,"").split(',')
                console.log(usersList.value)
            }).catch(e => {
                console.log(e)
                usersList.value = []
            })
        }

        const getUserTweets = (user) => {
            userApi.getTweetsData(user).then(data => {
                usersData.value.push(data)
                console.log(data)
            }).catch(err => {
                console.log(err)
            })
        }

        onMounted(() => {
          getUserList()
        })
        watch(usersList, () => {
            usersData.value = []
            for (let i = 0; i < usersList.value.length; i++) {
                getUserTweets(usersList.value[i])
            }
        })

        return {
            usersList,
            usersData,
            getUserTweets,
            getUserList,
            onExit
        }
    }
}
</script>
<style lang="less" scoped>
  .container {
    font-size: 12px;
    color: #666;

    .w1210 {
      width: 820px;
      margin: 0 auto;
    }

    .bd {
      margin-top: 2px;
      overflow: hidden;
      .aside-wrap {
        float: left;
        width: 210px;
        margin-bottom: 20px;
        :deep(.el-button) {
          background-color: #f7f7f7;
          color: #606266;
          border: 1px solid #DCDFE6;
          border-radius: 0;
          &:hover {
            color: red;
          }
        }
      }
      .detail {
        float: right;
        width: 600px;
        min-height: 220px;
      }
    }
  }
</style>
