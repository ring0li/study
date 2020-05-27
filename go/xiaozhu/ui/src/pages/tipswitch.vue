<style lang="less" scoped>
@fontsize1: 16px;
@iconsize1: 36px;
.listt {
  font-size: @fontsize1;
}
</style>
<template>
  <div>
    <br />
    <br />
    <van-cell-group>
      <van-switch-cell
        v-model="switchval"
        class="listt"
        @input="onInput"
        title="关闭打卡推送提醒"
        icon="bullhorn-o"
      />
    </van-cell-group>
  </div>
</template>

<script>
import api from "./data";

export default {
  name: "tipswitch",
  components: {},
  data() {
    return {
      switchval: true,
      openid: this.$route.query.openid
    };
  },
  mounted() {
    document.title = "设置打卡提醒";

    this._getTipData();
  },
  methods: {
    onInput(checked) {
      // this.$notify({ type: 'success', message: '通知内容' });

      if (checked == false) {
        this.$dialog
          .confirm({
            title: "提示",
            message: "关掉后将不再接收所有打卡提醒了？"
          })
          .then(() => {
            this._setTipData();
          })
          .catch(() => {});
      } else {
        this._setTipData();
      }
    },
    _getTipData() {
      let _this = this;
      api.getTipData({ openId: this.openid }).then(res => {
        if (res.errorCode == 0) {
          _this.switchval = res.data.switchval;
        }
      });
    },
    _setTipData() {
      let _this = this;
      api
        .setTipData({ openId: this.openid, switchval: this.switchval })
        .then(res => {
          if (res.errorCode == 0) {
            _this.$toast.success("设置成功");
          } else {
            _this.$toast(res.message);
          }
        });
    }
  }
};
</script>