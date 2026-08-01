<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import {
  listDdnsDomains,
  ddnsSelfIp,
  syncDdns,
  type DdnsDomainEntry,
  type DdnsSelfIp
} from '@/api/ddns'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

// Only records with console_enabled (synctype first digit '1') may be pushed
// from the console — FrontSyncHandler returns 403 otherwise. We fetch the
// full list once on mount and filter client-side.
const allRecords = ref<DdnsDomainEntry[]>([])
const selfIp = ref<DdnsSelfIp | null>(null)
const loading = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

// `key` is the composite subdomain.domain.PROTOCOL — stable & unique.
const form = reactive({ key: '', ip: '' })

// Stable, unique key matching the server's MetaMap key shape.
function rowKey(row: DdnsDomainEntry): string {
  return `${row.subdomain}.${row.domain}.${row.protocol}`
}

const consoleRecords = computed(() => allRecords.value.filter((r) => r.console_enabled))

const options = computed(() =>
  consoleRecords.value.map((r) => ({
    label: `${r.name} (${r.protocol === 'IPV6' ? t('ddnsView.tagProtocolV6') : t('ddnsView.tagProtocolV4')})`,
    value: rowKey(r)
  }))
)

const selected = computed(
  () => consoleRecords.value.find((r) => rowKey(r) === form.key) ?? null
)

// The caller's current IP, restricted to the family matching the selected
// record's protocol (a record can only hold one family).
const currentIp = computed(() => {
  const s = selfIp.value
  if (!s || !selected.value) return ''
  return selected.value.protocol === 'IPV6' ? (s.ipv6 ?? '') : (s.ipv4 ?? '')
})

const selfIpHasAny = computed(() => {
  const s = selfIp.value
  return !!(s && (s.ipv4 || s.ipv6))
})

// If we detected an IP but it's the wrong family for the selected record,
// surface a hint instead of silently disabling the link.
const familyMismatch = computed(
  () => selfIpHasAny.value && selected.value !== null && !currentIp.value
)

// Light client-side IP check; the server is the final authority.
const IPV4 = /^(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}$/
const IPV6 = /^[0-9a-fA-F:]+$/
function isIp(v: string): boolean {
  const s = (v || '').trim()
  if (!s) return false
  if (IPV4.test(s)) return true
  return s.includes(':') && IPV6.test(s)
}

// Reactive rules so validation messages follow the active locale.
const rules = computed<FormRules>(() => ({
  key: [{ required: true, message: t('ddnsUpdate.ruleRecord'), trigger: 'change' }],
  ip: [
    { required: true, message: t('ddnsUpdate.ruleIp'), trigger: 'blur' },
    {
      validator: (_r, value, cb) => {
        if (!value || !isIp(value)) cb(new Error(t('ddnsUpdate.ruleIpFormat')))
        else cb()
      },
      trigger: 'blur'
    }
  ]
}))

function useCurrentIp() {
  if (currentIp.value) form.ip = currentIp.value
}

async function loadData() {
  loading.value = true
  try {
    const [records, ip] = await Promise.all([listDdnsDomains(), ddnsSelfIp()])
    allRecords.value = records
    selfIp.value = ip

    // Pre-select from the route query (navigated from the View page's 更新 button).
    const q = route.query
    if (q.domain && q.subdomain && q.protocol) {
      const hit = records.find(
        (r) =>
          r.domain === q.domain &&
          r.subdomain === q.subdomain &&
          r.protocol === q.protocol
      )
      if (hit && hit.console_enabled) form.key = rowKey(hit)
    }
    // Otherwise default to the first console-enabled record.
    if (!form.key && consoleRecords.value.length) {
      form.key = rowKey(consoleRecords.value[0])
    }
  } catch {
    // surfaced by interceptor
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  const rec = selected.value
  if (!rec) return

  // Client-side "same ip, skip" pre-check: the server's only idempotent case
  // returns message "same ip, skip" with data:null, which the axios
  // envelope-unwrapper discards (it only returns `data`). The server's same-ip
  // test is `*ipMeta.Ip == ip`, equivalent to comparing against rec.ip here —
  // so we short-circuit with an info message and avoid the round-trip.
  if (rec.ip && rec.ip === form.ip.trim()) {
    ElMessage.info(t('ddnsUpdate.sameIp'))
    return
  }

  submitting.value = true
  try {
    await syncDdns(rec.domain, rec.subdomain, rec.protocol, { ip: form.ip.trim() })
    ElMessage.success(t('ddnsUpdate.success'))
    router.replace('/ddns/view')
  } catch {
    // 403/400/500 messages surfaced by the global interceptor
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>

<template>
  <div v-loading="loading" class="ddns-update-page">
    <el-card shadow="never">
      <template #header>
        <span class="card-title">{{ t('ddnsUpdate.title') }}</span>
      </template>

      <el-alert
        v-if="!loading && consoleRecords.length === 0"
        type="info"
        :closable="false"
        show-icon
        :title="t('ddnsUpdate.alertTitle')"
        :description="t('ddnsUpdate.alertDesc')"
      />

      <el-form
        v-else
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        style="max-width: 640px"
      >
        <el-form-item :label="t('ddnsUpdate.record')" prop="key">
          <el-select
            v-model="form.key"
            :placeholder="t('ddnsUpdate.placeholderRecord')"
            style="width: 100%"
          >
            <el-option v-for="o in options" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>

        <el-form-item :label="t('ddnsUpdate.currentIp')">
          <div class="current-ip-row">
            <span class="current-ip-value">
              <template v-if="!selfIpHasAny">{{ t('ddnsUpdate.currentIpMissing') }}</template>
              <template v-else-if="currentIp">{{ currentIp }}</template>
              <template v-else-if="familyMismatch">
                {{ t('ddnsUpdate.currentIpFamilyMismatch', { protocol: selected?.protocol }) }}
              </template>
              <template v-else>-</template>
            </span>
            <el-button
              text
              type="primary"
              :disabled="!currentIp"
              @click="useCurrentIp"
            >
              {{ t('ddnsUpdate.useCurrentIp') }}
            </el-button>
          </div>
        </el-form-item>

        <el-form-item :label="t('ddnsUpdate.ip')" prop="ip">
          <el-input v-model="form.ip" :placeholder="t('ddnsUpdate.placeholderIp')" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            {{ t('ddnsUpdate.submit') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.ddns-update-page {
  max-width: 960px;
}

.card-title {
  font-weight: 600;
}

.current-ip-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.current-ip-value {
  color: #606266;
  font-size: 14px;
}
</style>
