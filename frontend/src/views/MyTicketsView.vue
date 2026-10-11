<template>
  <section class="tickets-page">
    <div class="tickets-container">
      <header class="tickets-header">
        <span class="eyebrow">CINESTREAM</span>
        <h1>Vé của tôi</h1>
        <p>Quản lý và xem lại các vé bạn đã đặt.</p>
      </header>

      <div v-if="loading" class="state-message">Đang tải danh sách vé…</div>
      <div v-else-if="error" class="state-message error-message">
        {{ error }}
        <button class="retry-button" @click="loadTickets">Thử lại</button>
      </div>
      <div v-else-if="!tickets.length" class="state-message empty-message">
        <TicketCheck :size="34" />
        <strong>Bạn chưa có vé nào</strong>
        <span>Chọn phim và đặt suất chiếu để vé xuất hiện tại đây.</span>
        <router-link to="/booking" class="browse-button">Đặt vé xem phim</router-link>
      </div>

      <div v-else class="ticket-section">
        <nav class="ticket-tabs" aria-label="Lọc vé theo trạng thái">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            class="ticket-tab"
            :class="{ active: selectedTab === tab.id }"
            :aria-pressed="selectedTab === tab.id"
            @click="selectedTab = tab.id"
          >{{ tab.label }} <span>{{ tabTickets(tab.id).length }}</span></button>
        </nav>
        <div v-if="!visibleTickets.length" class="state-message empty-message filtered-empty">
          <strong>{{ emptyTabMessage }}</strong>
          <span>Chọn trạng thái khác để xem các vé đã đặt.</span>
        </div>
        <div v-else class="ticket-list">
        <article v-for="ticket in visibleTickets" :key="ticket.code" class="ticket-card">
          <div class="ticket-main">
            <div class="ticket-heading">
              <div>
                <span class="ticket-kicker">MÃ VÉ · {{ ticket.code }}</span>
                <h2>{{ ticket.showtime.movie_title }}</h2>
              </div>
              <span class="ticket-status" :class="ticket.status.toLowerCase()">
                {{ statusLabel(ticket.status) }}
              </span>
            </div>
            <div class="ticket-details">
              <div><span>Ngày chiếu</span><strong>{{ formatDate(ticket.showtime.start_at) }}</strong></div>
              <div><span>Suất chiếu</span><strong>{{ formatTime(ticket.showtime.start_at) }}</strong></div>
              <div><span>Rạp · Phòng</span><strong>{{ ticket.showtime.cinema?.name }} · {{ ticket.showtime.auditorium?.name }}</strong></div>
              <div><span>Ghế</span><strong>{{ ticket.seat_codes?.join(', ') || '—' }}</strong></div>
            </div>
          </div>
          <footer class="ticket-footer">
            <span>{{ ticket.payment_status === 'PAID' ? 'Đã thanh toán' : 'Tổng tiền' }}</span>
            <strong>{{ formatCurrency(ticket.total_amount) }}</strong>
          </footer>
        </article>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { TicketCheck } from '@lucide/vue';
import { fetchMyBookings, verifyBookingPayment } from '../services/booking.js';

const tickets = ref([]);
const selectedTab = ref('all');
const tabs = [
  { id: 'all', label: 'Tất cả' },
  { id: 'purchased', label: 'Đã mua' },
  { id: 'pending', label: 'Chờ thanh toán' },
  { id: 'cancelled', label: 'Đã hủy' },
  { id: 'expired', label: 'Đã hết hạn' },
];
const visibleTickets = computed(() => tabTickets(selectedTab.value));
const emptyTabMessage = computed(() => ({
  purchased: 'Bạn chưa có vé đã mua',
  pending: 'Không có vé đang chờ thanh toán',
  cancelled: 'Bạn chưa có vé đã hủy',
  expired: 'Bạn chưa có vé hết hạn',
}[selectedTab.value] || 'Bạn chưa có vé nào ở trạng thái này'));
const loading = ref(true);
const error = ref('');

function tabTickets(tab) {
  if (tab === 'all') return tickets.value;
  const statusByTab = { purchased: 'CONFIRMED', pending: 'AWAITING_PAYMENT', cancelled: 'CANCELLED', expired: 'EXPIRED' };
  return tickets.value.filter(ticket => ticket.status === statusByTab[tab]);
}

async function loadTickets() {
  loading.value = true;
  error.value = '';
  try {
    const history = await fetchMyBookings();
    tickets.value = await Promise.all(history.map(async (ticket) => {
      if (ticket.status !== 'AWAITING_PAYMENT') return ticket;
      try {
        const result = await verifyBookingPayment(ticket.code);
        return result.booking || ticket;
      } catch {
        return ticket;
      }
    }));
  } catch (err) {
    error.value = err.message || 'Không thể tải danh sách vé.';
  } finally {
    loading.value = false;
  }
}

function formatDate(value) {
  return new Intl.DateTimeFormat('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric', timeZone: 'Asia/Ho_Chi_Minh' }).format(new Date(value));
}

function formatTime(value) {
  return new Intl.DateTimeFormat('vi-VN', { hour: '2-digit', minute: '2-digit', timeZone: 'Asia/Ho_Chi_Minh' }).format(new Date(value));
}

function formatCurrency(amount) {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND', maximumFractionDigits: 0 }).format(amount || 0);
}

function statusLabel(status) {
  return ({ CONFIRMED: 'Đã mua', AWAITING_PAYMENT: 'Chờ thanh toán', EXPIRED: 'Đã hết hạn', CANCELLED: 'Đã hủy' })[status] || status;
}

onMounted(loadTickets);
</script>

<style scoped>
.tickets-page { min-height: 80vh; padding: 116px 24px 64px; color: var(--text-primary); }
.tickets-container { width: min(100%, 920px); margin: 0 auto; }
.tickets-header { margin-bottom: 28px; }
.eyebrow, .ticket-kicker { color: var(--accent-cyan); font-size: .75rem; font-weight: 800; letter-spacing: .15em; }
.tickets-header h1 { margin: 8px 0; font-family: var(--font-heading); font-size: clamp(2rem, 5vw, 3rem); }
.tickets-header p, .ticket-details span, .ticket-footer span { color: var(--text-secondary); }
.state-message { display: flex; align-items: center; justify-content: center; gap: 14px; min-height: 180px; padding: 24px; border: 1px solid var(--border-subtle); border-radius: 16px; background: var(--bg-card); color: var(--text-secondary); }
.empty-message { flex-direction: column; text-align: center; }
.empty-message strong { color: var(--text-primary); font-size: 1.1rem; }
.ticket-list { display: grid; gap: 16px; }
.ticket-tabs { display: flex; gap: 8px; margin-bottom: 18px; padding-bottom: 4px; overflow-x: auto; }
.ticket-tab { display: inline-flex; align-items: center; gap: 8px; flex: none; padding: 10px 14px; border: 1px solid var(--border-subtle); border-radius: 999px; color: var(--text-secondary); background: var(--bg-card); font: inherit; font-size: .88rem; cursor: pointer; transition: .2s ease; }
.ticket-tab span { min-width: 22px; padding: 2px 6px; border-radius: 999px; background: rgba(255,255,255,.08); text-align: center; font-size: .75rem; }
.ticket-tab.active { border-color: var(--accent-cyan); color: var(--text-primary); background: rgba(0,242,254,.1); }
.ticket-tab.active span { color: #041014; background: var(--accent-cyan); }
.filtered-empty { min-height: 140px; }
.ticket-card { overflow: hidden; border: 1px solid var(--border-subtle); border-radius: 16px; background: var(--bg-card); }
.ticket-main { padding: 22px 24px; }
.ticket-heading { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; }
.ticket-heading h2 { margin: 7px 0 0; font-size: 1.25rem; }
.ticket-status { flex: none; padding: 7px 11px; border-radius: 999px; background: rgba(148,163,184,.14); color: #cbd5e1; font-size: .8rem; font-weight: 700; }
.ticket-status.confirmed { color: #86efac; background: rgba(34,197,94,.13); }
.ticket-status.awaiting_payment { color: #fcd34d; background: rgba(245,158,11,.13); }
.ticket-status.cancelled, .ticket-status.expired { color: #fca5a5; background: rgba(239,68,68,.13); }
.ticket-details { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 16px 24px; margin-top: 22px; }
.ticket-details div { display: grid; gap: 5px; }
.ticket-details span, .ticket-footer { font-size: .85rem; }
.ticket-footer { display: flex; justify-content: space-between; align-items: center; padding: 14px 24px; border-top: 1px dashed var(--border-subtle); background: rgba(255,255,255,.025); }
.ticket-footer strong { color: var(--text-primary); font-size: 1.05rem; }
.browse-button, .retry-button { padding: 10px 16px; border: 0; border-radius: 10px; background: var(--accent-primary, #6d5dfc); color: white; font-weight: 700; text-decoration: none; cursor: pointer; }
.error-message { flex-wrap: wrap; color: #fca5a5; }
@media (max-width: 600px) { .tickets-page { padding: 100px 14px 44px; } .ticket-main { padding: 18px 16px; } .ticket-heading { flex-direction: column; } .ticket-details { grid-template-columns: 1fr; gap: 13px; } .ticket-footer { padding: 14px 16px; } }
</style>
