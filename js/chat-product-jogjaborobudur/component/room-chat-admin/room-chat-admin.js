$(document).ready(function () {
    let chatList = [];
    let adminChannel = null;

    function formatDateTime(datetime) {
        if (!datetime) return "";
        const d = new Date(datetime);
        return d.toLocaleString("id-ID", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit"
        });
    }

    function initPusher() {
        if (!window.pusher) {
            Pusher.logToConsole = false;
            window.pusher = new Pusher("4281def63ee8450bd642", {
                cluster: "ap1",
                forceTLS: true
            });
        }
    }

    function connectChatListAdminPusher() {
        initPusher();

        if (adminChannel) {
            adminChannel.unbind_all();
            window.pusher.unsubscribe("admin-sessions");
        }

        adminChannel = window.pusher.subscribe("admin-sessions");

        adminChannel.bind("session-update", function (msg) {
            const idx = chatList.findIndex(s => s.token === msg.token);

            const normalized = {
                fullname: msg.fullname,
                user_session: msg.user_session,
                token: msg.token,
                product_id: msg.product_id,
                product_name: msg.product_name,
                thumbnail: msg.thumbnail,
                updated_at: msg.updated_at,
                is_read_admin: msg.is_read_admin
            };

            if (idx !== -1) {
                chatList.splice(idx, 1);
            }

            chatList.unshift(normalized);
            drawChatList(chatList);
        });
    }

    window.renderChatListAdmin = function () {
        $("#chatList").html(`<p class="text-center">Loading...</p>`);

        $.get("https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/admin/chat-sessions")
            .then(res => {
                chatList = res.sessions;
                drawChatList(chatList);
            })
            .catch(() => {
                $("#chatList").html(`<p class="text-danger">Failed load chat</p>`);
            });
    };

    function drawChatList(list) {
        const container = $("#chatList");
        container.html("");

        if (!list.length) {
            container.html(`<p class="text-center text-muted">No chat yet</p>`);
            return;
        }

        list.forEach(item => {
            container.append(`
                <div class="chat-item ${item.is_read_admin ? "" : "unread"}"
                     data-token="${item.token}"
                     data-session="${item.user_session}"
                     data-product-id="${item.product_id}">
                    <img src="${item.thumbnail}">
                    <div class="chat-item-title">
                        <div class="chat-title-row">
                            ${item.product_name}
                            ${item.is_read_admin ? "" : `<span class="badge-new">New</span>`}
                        </div>
                        <p class="text-secondary text-capitalize" style="font-size:12px;">Fullname: ${item.fullname}</p>
                        <p class="text-secondary" style="font-size:12px;">
                            Last update: ${formatDateTime(item.updated_at)}
                        </p>
                    </div>
                </div>
                <hr>
            `);
        });
    }

    $(document).on("click", ".chat-item", function () {
        const token = $(this).data("token");
        const session = $(this).data("session")
        const productId = $(this).data("product-id")
        const prodid = Number(productId)
        
        $.ajax({
            url: `https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/open/user/${token}/admin`,
            method: "PATCH"
        }).always(() => {
            localStorage.setItem("activeChatToken", token);

            $("#chatRoomWidgetAdmin").fadeOut(150);
            $("#chatProductWidgetAdmin").fadeIn(150);

            window.loadChatRoom(token, session, prodid);
        });
    });

    $(document).on("click", "#chatCloseBtnRoomAdmin", function () {
        localStorage.removeItem("activeChatToken");
        $("#chatRoomWidgetAdmin").fadeOut(150);
        $("#chatLauncherAdmin").fadeIn(150);
        $("#chatLauncherAdminMobile").fadeIn(150);
        $("#chatOverlayMobileAdmin").fadeOut(150);
    });

    window.renderChatListAdmin();
    connectChatListAdminPusher();
});
