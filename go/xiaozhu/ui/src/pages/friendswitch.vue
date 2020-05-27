<style lang="less" scoped>
@fontsize1: 16px;
@fontsize2: 26px;
@iconsize1: 36px;
.listt {
  font-size: @fontsize1;
  padding-left: @fontsize1;
}
.licon {
  height: @iconsize1;
  width: @iconsize1;
}
.ldex {
  font-size: @fontsize2;
  padding-right: 5px;
  color: #eae8e8;
}
.content {
  text-align: center;
  padding-top: 5px;
}
</style>
<template>
  <div>
    <van-cell-group>
      <van-cell v-for="(item,idx) in flist" :key="idx" :center="true">
        <template slot="icon">
          <van-image class="licon" :src="item.HeadImgUrl" :round="true" cover />
        </template>
        <template slot="title">
          <div class="listt">{{item.friendName}}</div>
        </template>
        <template slot="default">
          <div>
            <van-icon name="close" class="ldex" @click="_delaction(item)" />
            <van-switch v-model="item.switchval" @input="onInput(item)" size="24px" />
          </div>
        </template>
      </van-cell>
    </van-cell-group>
    <div class="content">
      <van-button icon="friends-o" size="normal" @click="addfriendtip">加好友</van-button>
    </div>
  </div>
</template>

<script>
import api from "./data";

export default {
  name: "friendswitch",
  components: {},
  data() {
    return {
      openid: this.$route.query.openid,
      flist: []
    };
  },
  mounted() {
    document.title = "设置好友提醒";
    this._getFriendTipData();
  },
  methods: {
    addfriendtip() {
      this.$dialog
        .alert({
          message: "把打卡图分享给好友，好友扫码后，关联成好友"
        })
        .then(() => {
          // on close
        });
    },
    onInput(item) {
      if (item.switchval == false) {
        this.$dialog
          .confirm({
            title: "提示",
            message: "关掉后将不再接收好友提醒了？"
          })
          .then(() => {
            this._setFriendTipData(item);
          })
          .catch(() => {
            item.switchval = true;
          });
      } else {
        this._setFriendTipData(item);
      }
    },
    _delaction(item) {
      let _this = this;
      this.$dialog
        .confirm({
          title: "提示",
          message: "要解除好友的关系吗？"
        })
        .then(() => {
          _this._delFriend(item);
        })
        .catch(() => {});
    },
    _getFriendTipData() {
      let _this = this;
      api.getFriendTipData({ openId: this.openid }).then(res => {
        if (res.errorCode == 0) {
          _this.flist = res.data;
        }
      });
    },
    _setFriendTipData(item) {
      let _this = this;
      let _f = Object.assign({ openId: this.openid }, item);
      api.setFriendTipData(_f).then(res => {
        if (res.errorCode == 0) {
          _this.$toast.success("设置成功");
        } else {
          _this.$toast(res.message);
        }
      });
    },
    _delFriend(item) {
      let _this = this;
      api.delFriend(item).then(res => {
        if (res.errorCode == 0) {
          _this._getFriendTipData();
          _this.$toast.success("解除成功");
        } else {
          _this.$toast(res.message);
        }
      });
    }
  }
};
</script>

<style lang="less" scoped>
</style>
