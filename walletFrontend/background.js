
// 监听来自内容脚本的消息
// chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
//     // 弹出钱包页面
//     if (request.action === "openPopup") {
//         chrome.windows.create({
//             url: "/views/sign.html",
//             type: "popup",
//             width: 360,
//             height: 600,
//             top: 0,
//             left: screen.availWidth - 180
//         });
//     }

// });