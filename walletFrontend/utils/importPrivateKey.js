$(document).ready(function () {
    chrome.storage.local.get('wallet', function (result) {
        // 获取存储的值
        let account = result.wallet;

        // 在这里使用获取到的值
        let addr = account.blockchain_address;
        if (addr != undefined) {
            $('#address').html(formatStr(addr));
            $("#accountValue").html(account.amount)
        }
        // 执行其他操作
    });


    $("#accountDropdown").click(() => {
        chrome.storage.local.get('walletList', function (result) {
            let content = ""
            for (let i = 0; i < result.walletList.length; i++) {
                content += `<li><a class="dropdown-item transAccount" href="#" value="${result.walletList[i].blockchain_address}">${formatStr(result.walletList[i].blockchain_address)}</a></li>`
            }
            $("#accountList").html(content)
        });
    })

    $("#copy").click(() => {
        chrome.storage.local.get('wallet', function (result) {
            let account = result.wallet;
            let addr = account.blockchain_address;
            copyToClipboard(addr);

        });
    })

    $(document).on("click",".record ul li a",function(){
        let hash = $(this).attr("hash")
        copyToClipboard(hash)
    })

    // 获取元素
    var importMethodSelect = $('#importMethod');
    var privateKeyContainer = $('#privateKeyContainer');
    var addContainer = $('#addContainer');

    // 隐藏初始状态的容器
    addContainer.hide();

    // 添加事件监听器
    importMethodSelect.on('change', function () {
        if (importMethodSelect.val() === 'privateKey') {
            privateKeyContainer.show();
            addContainer.hide();
        } else if (importMethodSelect.val() === 'addAccount') {
            privateKeyContainer.hide();
            addContainer.show();
        }
    });

    $(document).on("click", "#accountList li a", function () {
        let address = $(this).attr("value")
        $('#address').html(formatStr(address));
        let data ={
            blockchain_address:address
        }
        $.ajax({
            url:"http://127.0.0.1:8080/amount",
            type:"post",
            data:JSON.stringify(data),
            success:function(response){
                if(response.message=="success"){
                    $("#accountValue").html(response.amount)
                    chrome.storage.local.get('walletList', function (result) {
                        for (let i = 0; i < result.walletList.length; i++) {
                            console.log(result.walletList[i].blockchain_address === address);
                            if (result.walletList[i].blockchain_address === address) {
                                console.log(result.walletList[i]);
                                let wallet = {
                                    public_key: result.walletList[i].public_key,
                                    private_key: result.walletList[i].private_key,
                                    blockchain_address: result.walletList[i].blockchain_address,
                                    amount:response.amount
                                }
                                
                                chrome.storage.local.set({ wallet: wallet }, function () {
                                });
                            }
                        }
                    });
                }
            }
        })
       
    })

    $('#importPrivateKeyBtn').click(function () {
        let privateKey = $('#privateKey').val()
        $.ajax({
            url: 'http://127.0.0.1:8080/loadWallet',
            type: 'POST',
            data: {
                privateKey: privateKey
            },
            success: function (response) {
                response = JSON.parse(response)
                let wallet = {
                    public_key: response["public_key"],
                    private_key: response["private_key"],
                    blockchain_address: response["blockchain_address"],
                    amount:0
                }
                chrome.storage.local.get("walletList", function (result) {
                    if (result.walletList == undefined) {
                        let wl = []
                        wl.push(wallet)
                        chrome.storage.local.set({ walletList: wl }, function () {
                            console.log("账户列表保存成功")
                        })
                    } else {
                        let wl = result.walletList
                        wl.push(wallet)
                        chrome.storage.local.set({ walletList: wl }, function () {
                            console.log("账户列表保存成功")
                        })
                    }
                })

                chrome.storage.local.set({ wallet: wallet }, function () {
                    console.log('数据已保存');
                });
                $('#address').html(formatStr(response["blockchain_address"]));
            }
        })
    })
    $('#addAccountBtn').click(function (e) {
        $.ajax({
            url: "http://127.0.0.1:8080/wallet",
            type: "POST",
            success: function (response) {
                response = JSON.parse(response)

                let wallet = {
                    public_key: response["public_key"],
                    private_key: response["private_key"],
                    blockchain_address: response["blockchain_address"],
                    amount:0
                }

                chrome.storage.local.get("walletList", function (result) {
                    if (result.walletList == undefined) {
                        let wl = []
                        wl.push(wallet)
                        chrome.storage.local.set({ walletList: wl }, function () {
                            console.log("账户列表保存成功")
                        })
                    } else {
                        let wl = result.walletList
                        wl.push(wallet)
                        chrome.storage.local.set({ walletList: wl }, function () {
                            console.log("账户列表保存成功")
                        })
                    }
                })

                $("#inputPublic").val(response["public_key"]);
                $("#inputPrivateKey").val(response["private_key"]);
                $("#inputAddress").val(response["blockchain_address"]);
                chrome.storage.local.set({ wallet: wallet }, function () {
                    console.log('数据已保存');
                });
                $('#address').html(formatStr(response["blockchain_address"]));
            },
            error: function (error) {
                console.error(error);
            },
        });
    })


})

$("#buttonSubmit").click(function(){
    chrome.storage.local.get('wallet', function (result) {
        let account = result.wallet;
        let FromAddr = account.blockchain_address;
        let privKey = account.private_key
        let pubKey = account.public_key
        let ToAddr = $("#inputReceiveAddress").val()
        let value = $("#inputAmount").val()
        let data ={
            sender_public_key:pubKey,
            sender_blockchain_address:FromAddr,
            recipient_blockchain_address:ToAddr,
            sender_private_key:privKey,
            value:value
        }
        $.ajax({
            url: "http://127.0.0.1:8080/createTransaction",
            type:"post",
            data:JSON.stringify(data),
            success:function(response){
                if (response.message=="success") {
                    alert("交易成功")
                    window.location.reload()
                }else{
                    alert("交易失败")
                }
            }
        })


    });
    

})

$("#nav-contact-tab").click(function(){
    $.ajax({
        url: "http://127.0.0.1:5000/transactionRecords",
        type:"get",
        success:function(response){
            let content = ""
            for (let i = 0; i < response.transactions.length; i++) {
                content += `<div class="row record" >
                <ul class="list-group">
                    <li class="list-group-item">From： <span>${response.transactions[i].from}</span></li>
                    <li class="list-group-item">To： <span>${response.transactions[i].to}</span></li>
                    <li class="list-group-item">Value： <span>${response.transactions[i].value}</span></li>
                    <li class="list-group-item">Hash： <span>${formatStr(response.transactions[i].hash)}</span></li>
                    <li class="list-group-item text-center"><a href="javascript:void(0)" class="copyHash" style="text-decoration: none;" hash="${response.transactions[i].hash}">复制 Hash</a></li>
                  </ul>
            </div>`
            }
            $("#nav-contact").html(content)
        }
    })
})

// document.getElementById("cancle").addEventListener("click", function () {
//     chrome.windows.getCurrent(function (window) {
//         chrome.windows.remove(window.id);
//     });
// });


// document.getElementById("confirm").addEventListener("click", function () {

//     $.ajax({
//         url: 'http://127.0.0.1:5000/getSignature',
//         type: 'Get',
//         success: function (response) {
//             if (response.add_status == "add_err") {
//                 alert(response.add_msg)
//             } else {
//                 alert("添加成功")
//             }
//         },

//     });
//     chrome.windows.remove(window.id);
// });


function formatStr(addr) {
    const prefix = addr.substring(0, 6);
    const suffix = addr.substring(addr.length - 6);
    const result = `${prefix}...${suffix}`;
    return result
}

function copyToClipboard(text) {
    var dummyElement = $('<textarea>').val(text).appendTo('body').select();
    document.execCommand('copy');
    dummyElement.remove();
}



// function getAccounts(accounts, callback) {
//     chrome.storage.local.get(accounts, function(result) {
//         let storedAccount = result.accounts;
//         if (!storedAccount) {
//             callback(storedAccount);
//         } else {
//             callback(null);
//         }
//     });
// }



