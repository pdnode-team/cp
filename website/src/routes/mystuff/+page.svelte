<script lang="ts">
    import pb from '$lib/pocketbase';
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';

    type ViewMode = 'cps' | 'characters';
    let mode = $state<ViewMode>('cps');
    let items = $state<any[]>([]);
    let isLoading = $state(true);

    async function fetchMyData(currentMode: ViewMode) {
        // 安全守卫：如果没有登录，直接跳走
        if (!pb.authStore.isValid || !pb.authStore.record) {
            goto('/login');
            return;
        }

        isLoading = true;
        const userId = pb.authStore.record.id;

        try {
            const collection = currentMode === 'cps' ? 'cps' : 'characters';
            const options = {
                // 关键点：只获取 owner 是当前用户的记录
                filter: `owner = "${userId}"`,
                sort: '-created',
                expand: currentMode === 'cps' ? 'characters' : ''
            };

            items = await pb.collection(collection).getFullList(options);
        } catch (err) {
            console.error("Fetch error:", err);
        } finally {
            isLoading = false;
        }
    }

    // 监听 mode 变化自动刷新
    $effect(() => {
        fetchMyData(mode);
    });
</script>

<div class="mx-auto max-w-6xl px-6 py-16">
    <div class="mb-12">
        <h1 class="text-4xl font-bold md:text-5xl">My Stuff</h1>
        <p class="text-base-content/60 mt-2">Manage the CPs and Characters you've created.</p>
        
        <div role="tablist" class="tabs tabs-boxed mt-6 w-fit">
            <button 
                role="tab" 
                class="tab {mode === 'cps' ? 'tab-active' : ''}" 
                onclick={() => mode = 'cps'}
            >My CPs</button>
            <button 
                role="tab" 
                class="tab {mode === 'characters' ? 'tab-active' : ''}" 
                onclick={() => mode = 'characters'}
            >My Characters</button>
        </div>
    </div>

    {#if isLoading}
        <div class="flex justify-center py-20">
            <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>
    {:else if items.length === 0}
        <div class="card bg-base-200 py-20 text-center border-2 border-dashed border-base-300">
            <div class="card-body items-center">
                <h2 class="card-title text-2xl opacity-40">You haven't created any {mode} yet.</h2>
                <div class="card-actions mt-4">
                    <a href="/create" class="btn btn-primary">Create Your First {mode === 'cps' ? 'CP' : 'Char'}</a>
                </div>
            </div>
        </div>
    {:else}
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
            {#each items as item}
                <div class="card bg-base-100 shadow-xl border border-base-200">
                    <div class="card-body p-5">
                        <h2 class="card-title">{item.name}</h2>
                        <div class="card-actions justify-end mt-4">
                            <a href="/edit/{mode}/{item.id}" class="btn btn-sm btn-ghost">Edit</a>
                            <a href="/{mode}/{item.id}" class="btn btn-sm btn-primary">View</a>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>