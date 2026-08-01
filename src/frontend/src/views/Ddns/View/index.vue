<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { View as ViewIcon, Refresh, EditPen } from '@element-plus/icons-vue'
import {
  listDdnsDomains,
  ddnsInfo,
  type DdnsDomainEntry,
  type DdnsHistoryEntry,
  type DdnsInfo
} from '@/api/ddns'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()

// View-only enrichment on top of the /domain record. `dip` / `consistent` /
// `history` only come from /info (which does a live DNS query), so we fetch
// /info per row after the list loads. `infoStatus` drives the loading/failed
// rendering of the dip & consistent cells.
type InfoStatus = 'idle' | 'loading' | 'ok' | 'error'

interface DdnsRow extends DdnsDomainEntry {
  dip?: string
  consistent?: boolean
  history?: DdnsHistoryEntry[]
  infoStatus: InfoStatus
}

const list = ref<DdnsRow[]>([])
const loading = ref(false)

// Detail dialog state.
const detailVisible = ref(false)
const detailLoading = ref(false)
const detailRow = ref<DdnsRow | null>(null)

// Stable, unique key matching the server's MetaMap key shape.
function rowKey(row: DdnsDomainEntry): string {
  return `${row.subdomain}.${row.domain}.${row.protocol}`
}

// Stable colors per provider — purely cosmetic.
function providerTagType(provider: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'primary' | 'danger'> = {
    TENCENT: 'success',
    ALIYUN: 'warning',
    GODADDY: 'info',
    GOOGLE: 'primary',
    CLOUDFLARE: 'danger',
    NONE: 'info'
  }
  return map[provider] ?? 'info'
}

// We render "未知/Unknown" (not "不一致/Inconsistent") whenever the DNS query
// hasn't resolved or failed: an empty dip is NOT a real mismatch (the server
// sets consistent=false simply because dip==='').
function showConsistent(row: DdnsRow): boolean {
  return row.infoStatus === 'ok' && !!row.dip
}

async function fetchList() {
  loading.value = true
  try {
    const rows = await listDdnsDomains()
    list.value = rows.map((r) => ({ ...r, infoStatus: 'idle' as InfoStatus }))
  } catch {
    // surfaced by interceptor
  } finally {
    loading.value = false
  }
  // Enrich every row with dip/consistent/history (live DNS queries). We fire
  // them all in parallel and tolerate per-row failures so one bad record
  // can't blank out the rest.
  await loadAllInfo()
}

async function loadAllInfo() {
  const targets = list.value.slice()
  if (targets.length === 0) return
  targets.forEach((row) => (row.infoStatus = 'loading'))
  const results = await Promise.allSettled(
    targets.map((r) => ddnsInfo(r.domain, r.subdomain, r.protocol))
  )
  targets.forEach((row, i) => {
    const res = results[i]
    if (res.status === 'fulfilled') {
      const info: DdnsInfo = res.value
      row.dip = info.dip
      row.consistent = info.consistent
      row.history = info.history
      row.infoStatus = 'ok'
    } else {
      row.infoStatus = 'error' // dip '-', consistent 未知
    }
  })
}

function handleDetail(row: DdnsRow) {
  // Reuse the eagerly-fetched cache so the dialog opens instantly.
  detailRow.value = row
  detailVisible.value = true
}

// Optional in-dialog refresh: re-query /info and write the result back onto
// the row (the single source of truth shared with the table).
async function refreshDetail() {
  const row = detailRow.value
  if (!row) return
  detailLoading.value = true
  try {
    const info = await ddnsInfo(row.domain, row.subdomain, row.protocol)
    row.dip = info.dip
    row.consistent = info.consistent
    row.history = info.history
    row.infoStatus = 'ok'
  } catch {
    // surfaced by interceptor
  } finally {
    detailLoading.value = false
  }
}

function handleUpdate(row: DdnsRow) {
  router.push({
    path: '/ddns/update',
    query: { domain: row.domain, subdomain: row.subdomain, protocol: row.protocol }
  })
}

onMounted(fetchList)
</script>

<template>
  <div v-loading="loading" class="ddns-view-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">{{ t('ddnsView.title') }}</span>
          <el-button text type="primary" :icon="Refresh" @click="fetchList">
            {{ t('ddnsView.refresh') }}
          </el-button>
        </div>
      </template>

      <el-table :data="list" border stripe :empty-text="t('ddnsView.empty')">
        <el-table-column prop="name" :label="t('ddnsView.columnName')" min-width="200" />
        <el-table-column :label="t('ddnsView.columnProvider')" width="130">
          <template #default="{ row }">
            <el-tag :type="providerTagType(row.provider)" effect="plain">
              {{ row.provider || '-' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('ddnsView.columnProtocol')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.protocol === 'IPV6' ? 'primary' : 'success'" effect="dark">
              {{ row.protocol === 'IPV6' ? t('ddnsView.tagProtocolV6') : t('ddnsView.tagProtocolV4') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('ddnsView.columnConsoleEnabled')" width="120">
          <template #default="{ row }">
            <el-tag :type="row.console_enabled ? 'success' : 'info'" effect="plain">
              {{ row.console_enabled ? t('ddnsView.tagConsoleEnabled') : t('ddnsView.tagConsoleDisabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('ddnsView.columnClientUpload')" width="110">
          <template #default="{ row }">
            <el-tag :type="row.client_upload ? 'success' : 'info'" effect="plain">
              {{ row.client_upload ? t('ddnsView.tagClientUploadYes') : t('ddnsView.tagClientUploadNo') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('ddnsView.columnIp')" min-width="160">
          <template #default="{ row }">{{ row.ip || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('ddnsView.columnUpdateTime')" min-width="170">
          <template #default="{ row }">{{ row.update_time || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('ddnsView.columnDip')" min-width="160">
          <template #default="{ row }">
            <span v-if="row.infoStatus === 'loading'" class="cell-loading">
              {{ t('ddnsView.loadingDns') }}
            </span>
            <span v-else>{{ row.dip || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('ddnsView.columnConsistent')" width="110">
          <template #default="{ row }">
            <el-tag v-if="!showConsistent(row)" type="info" effect="plain">
              {{ t('ddnsView.tagConsistentUnknown') }}
            </el-tag>
            <el-tag v-else-if="row.consistent" type="success" effect="plain">
              {{ t('ddnsView.tagConsistentYes') }}
            </el-tag>
            <el-tag v-else type="danger" effect="plain">
              {{ t('ddnsView.tagConsistentNo') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('ddnsView.columnOperation')" width="170" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="ViewIcon" @click="handleDetail(row)">
              {{ t('ddnsView.btnDetail') }}
            </el-button>
            <el-button
              v-if="row.console_enabled"
              size="small"
              type="primary"
              :icon="EditPen"
              @click="handleUpdate(row)"
            >
              {{ t('ddnsView.btnUpdate') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="detailVisible"
      :title="t('ddnsView.detail.title')"
      width="640px"
      destroy-on-close
    >
      <div v-loading="detailLoading">
        <template v-if="detailRow">
          <div class="detail-toolbar">
            <el-button text type="primary" :icon="Refresh" @click="refreshDetail">
              {{ t('ddnsView.refresh') }}
            </el-button>
          </div>
          <el-descriptions :column="1" border>
            <el-descriptions-item :label="t('ddnsView.detail.name')">
              {{ detailRow.name }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('ddnsView.detail.provider')">
              <el-tag :type="providerTagType(detailRow.provider)" effect="plain">
                {{ detailRow.provider || '-' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="t('ddnsView.detail.protocol')">
              <el-tag :type="detailRow.protocol === 'IPV6' ? 'primary' : 'success'" effect="dark">
                {{ detailRow.protocol === 'IPV6' ? t('ddnsView.tagProtocolV6') : t('ddnsView.tagProtocolV4') }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="t('ddnsView.detail.ip')">
              {{ detailRow.ip || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('ddnsView.detail.updateTime')">
              {{ detailRow.update_time || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('ddnsView.detail.dip')">
              {{ detailRow.dip || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('ddnsView.detail.consistent')">
              <el-tag v-if="!showConsistent(detailRow)" type="info" effect="plain">
                {{ t('ddnsView.tagConsistentUnknown') }}
              </el-tag>
              <el-tag v-else-if="detailRow.consistent" type="success" effect="plain">
                {{ t('ddnsView.tagConsistentYes') }}
              </el-tag>
              <el-tag v-else type="danger" effect="plain">
                {{ t('ddnsView.tagConsistentNo') }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="t('ddnsView.detail.consoleEnabled')">
              {{ detailRow.console_enabled ? t('ddnsView.tagConsoleEnabled') : t('ddnsView.tagConsoleDisabled') }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('ddnsView.detail.clientUpload')">
              {{ detailRow.client_upload ? t('ddnsView.tagClientUploadYes') : t('ddnsView.tagClientUploadNo') }}
            </el-descriptions-item>
          </el-descriptions>

          <h4 class="history-title">{{ t('ddnsView.detail.historyTitle') }}</h4>
          <el-table
            :data="detailRow.history || []"
            border
            size="small"
            :empty-text="t('ddnsView.detail.emptyHistory')"
          >
            <el-table-column prop="ip" :label="t('ddnsView.detail.columnHistoryIp')" min-width="180" />
            <el-table-column
              prop="time"
              :label="t('ddnsView.detail.columnHistoryTime')"
              min-width="180"
            />
          </el-table>
        </template>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.ddns-view-page {
  width: 100%;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-weight: 600;
}

.cell-loading {
  color: #909399;
  font-style: italic;
}

.detail-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 8px;
}

.history-title {
  margin: 20px 0 10px;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}
</style>
