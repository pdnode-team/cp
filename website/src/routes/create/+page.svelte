<script lang="ts">
	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';

	let errorText = $state('');
	let tags = $state<string[]>([]);
	let currentInput = $state('');
	let newCharDialog = $state() as HTMLDialogElement;

	let cpName = $state('');
	let cpDescription = $state('');
	let cpPictures: FileList | undefined = $state();
	let cpCharacter1 = $state('');
	let cpCharacter2 = $state('');

	const createCP = async () => {
		errorText = '';

		if (
			!cpName.trim() ||
			!cpDescription.trim() ||
			!cpPictures ||
			cpPictures.length < 1 ||
			cpPictures.length > 3 ||
			!cpCharacter1.trim() ||
			!cpCharacter2.trim()
		) {
			errorText = 'All fields must be filled out';
			return;
		}

		if (cpCharacter1 === cpCharacter2) {
			errorText = 'Character #1 and Character #2 cannot be the same';
			return;
		}

		const formData = new FormData();
		formData.append('name', cpName);
		formData.append('description', cpDescription);
		formData.append('owner', pb.authStore.record!.id);
		for (let file of cpPictures) {
			formData.append('images', file);
		}
		tags.forEach((tag) => {
			formData.append('tag_names', tag);
		});
		formData.append('characters', cpCharacter1);
		formData.append('characters', cpCharacter2);

		try {
			await pb.collection('cps').create(formData);
			window.location.pathname = '/explore';
		} catch (err: any) {
			errorText = err.data.data?.message ?? 'Create failed. Please try again.';

			const firstKey = Object.keys(err.data.data)[0];
			if (firstKey) {
				const friendlyMessage = `${firstKey}: ${err.data.data[firstKey].message}`;
				console.error(friendlyMessage);
				errorText = friendlyMessage;
				return;
			}
		}
	};

	function addTag(e: KeyboardEvent) {
		if (e.key === 'Enter' && currentInput.trim()) {
			e.preventDefault();
			if (tags.length >= 10) {
				errorText = 'Maximum of ten tags';
				return;
			}
			// 如果标签不存在则添加
			if (!tags.includes(currentInput.trim())) {
				tags.push(currentInput.trim());
			}
			currentInput = ''; // 清空输入
		}
	}

	const handleChangeSelect = (e: Event) => {
		const target = e.target as HTMLSelectElement;
		if (target.value === 'new') {
			target.value = '';
			newCharDialog.showModal();
			return;
		}
	};

	// New Char Dialog
	let newCharTags = $state<string[]>([]);
	let newCharTagsCurrentInput = $state('');
	let newCharErrorText = $state('');

	function newCharAddTag(e: KeyboardEvent) {
		if (e.key === 'Enter' && newCharTagsCurrentInput.trim()) {
			e.preventDefault();
			if (newCharTags.length >= 10) {
				newCharErrorText = 'Maximum of ten tags';
				return;
			}
			// 如果标签不存在则添加
			if (!newCharTags.includes(newCharTagsCurrentInput.trim())) {
				newCharTags.push(newCharTagsCurrentInput.trim());
			}
			newCharTagsCurrentInput = ''; // 清空输入
		}
	}

	let newCharName = $state('');
	let newCharDescription = $state('');
	let newCharOrigin = $state('');
	let newCharPictures: FileList | undefined = $state();

	const createCharacter = async () => {
		const formData = new FormData();
		newCharErrorText = '';

		if (!newCharName.trim() || !newCharDescription.trim() || !newCharOrigin.trim()) {
			newCharErrorText = 'All fields must be filled out';
			return;
		}

		if (newCharPictures && newCharPictures.length !== 0) {
			for (let file of newCharPictures) {
				formData.append('images', file);
			}
		}

		formData.append('name', newCharName);
		formData.append('description', newCharDescription);
		formData.append('origin', newCharOrigin);
		newCharTags.forEach((tag) => {
			formData.append('tag_names', tag);
		});

		formData.append('owner', pb.authStore.record!.id);

		try {
			await pb.collection('characters').create(formData);
			newCharDialog.close();
			newCharName = '';
			newCharDescription = '';
			newCharOrigin = '';
			newCharPictures = undefined;
			newCharTags = [];
			newCharTagsCurrentInput = '';
			newCharErrorText = '';
			reloadCharacters();
		} catch (err: any) {
			errorText = err.data.data?.message ?? 'Create failed. Please try again.';

			const firstKey = Object.keys(err.data.data)[0];
			if (firstKey) {
				const friendlyMessage = `${firstKey}: ${err.data.data[firstKey].message}`;
				console.error(friendlyMessage);
				newCharErrorText = friendlyMessage;
				return;
			}
		}
	};

	// Get Char(s)
	let characters = $state<any[]>([]);
	const reloadCharacters = async () => {
		characters = await pb.collection('characters').getFullList();
	};
	onMount(() => {
		if (!pb.authStore.isValid) {
			window.location.pathname = '/login';
			return
		}
		reloadCharacters();
	});
</script>

<!-- New Char Dialog -->
<dialog bind:this={newCharDialog} class="modal">
	<div class="modal-box flex flex-col gap-4">
		<h3 class="text-lg font-bold">Create a new character</h3>
		<div class="flex flex-col gap-4">
			<label class="form-control w-full">
				<div class="label"><span class="label-text font-medium">Your Character Name</span></div>
				<input
					type="text"
					name="characterName"
					placeholder="XXXX & YYYY"
					class="input-bordered input w-full"
					autocomplete="name"
					bind:value={newCharName}
					required
				/>
			</label>

			<label class="form-control w-full">
				<div class="label">
					<span class="label-text font-medium">Your Character Description</span>
				</div>
				<textarea
					name="characterpDescription"
					placeholder="description"
					class="textarea w-full"
					maxlength="1000"
					bind:value={newCharDescription}
					required
				></textarea>
			</label>

			<label class="form-control w-full">
				<div class="label"><span class="label-text font-medium">Your Character Origin</span></div>
				<input
					type="text"
					name="characterOrigin"
					placeholder="https://pdnode.com"
					class="input-bordered input w-full"
					autocomplete="name"
					bind:value={newCharOrigin}
				/>
			</label>

			<div class="form-control w-full max-w-xs">
				<label class="label" for="file-upload">
					<span class="label-text">Your Character Pictures (Option)</span>
				</label>

				<input
					type="file"
					id="file-upload"
					class="file-input-bordered file-input w-full file-input-primary"
					accept="image/*"
					multiple
					bind:files={newCharPictures}
				/>
			</div>

			<div class="flex w-full max-w-sm flex-col">
				<div class="label">
					<span class="label-text font-medium">Your Character Tags (Option)</span>
				</div>

				<input
					type="text"
					placeholder="Enter the tags and press Enter...."
					class="input-bordered input w-full"
					bind:value={newCharTagsCurrentInput}
					onkeydown={newCharAddTag}
				/>

				<div class="mt-2 flex flex-wrap gap-2">
					{#each newCharTags as tag, i}
						<div class="badge gap-2 badge-soft p-3 badge-primary">
							{tag}
							<button type="button" onclick={() => newCharTags.splice(i, 1)} class="text-xs"
								>✕</button
							>
						</div>
					{/each}
				</div>
			</div>

			<p class="text-red-500">{newCharErrorText}</p>

			<div class="form-control mt-6">
				<button type="submit" class="btn w-full text-lg btn-primary" onclick={createCharacter}
					>Create</button
				>
			</div>
			<div class="modal-action">
				<form method="dialog">
					<!-- if there is a button in form, it will close the modal -->
					<button class="btn">Close</button>
				</form>
			</div>
		</div>
	</div>
</dialog>
<!-- End of Dialog -->

<div class="flex min-h-[70vh] flex-col items-center justify-center px-4 py-12">
	<div class="card w-full max-w-md bg-base-100">
		<div class="card-body p-8">
			<h2 class="mb-2 text-center text-3xl font-bold">Create a CP</h2>
			<p class="mb-6 text-center text-base-content/60">
				All fields must be filled out unless otherwise specified.
			</p>
			<div class="flex flex-col gap-4">
				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Your CP Name</span></div>
					<input
						type="text"
						name="cpName"
						placeholder="XXXX & YYYY"
						class="input-bordered input w-full"
						autocomplete="name"
						bind:value={cpName}
						required
					/>
				</label>

				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Your CP Description</span></div>
					<textarea
						name="cpDescription"
						placeholder="description"
						class="textarea w-full"
						maxlength="1000"
						bind:value={cpDescription}
						required
					></textarea>
				</label>

				<div class="form-control w-full max-w-xs">
					<label class="label" for="file-upload">
						<span class="label-text">Your CP Pictures</span>
					</label>

					<input
						type="file"
						id="file-upload"
						class="file-input-bordered file-input w-full file-input-primary"
						accept="image/*"
						bind:files={cpPictures}
						multiple
						required
					/>
				</div>

				<div class="flex w-full max-w-sm flex-col">
					<div class="label"><span class="label-text font-medium">Your CP Tags</span></div>

					<input
						type="text"
						placeholder="Enter the tags and press Enter...."
						class="input-bordered input w-full"
						bind:value={currentInput}
						onkeydown={addTag}
					/>

					<div class="mt-2 flex flex-wrap gap-2">
						{#each tags as tag, i}
							<div class="badge gap-2 badge-soft p-3 badge-primary">
								{tag}
								<button type="button" onclick={() => tags.splice(i, 1)} class="text-xs">✕</button>
							</div>
						{/each}
					</div>
				</div>

				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Your Character #1</span></div>
					<select
						name="cpCharacter1"
						placeholder="Select a character"
						class="select-bordered select w-full"
						autocomplete="name"
						onchange={handleChangeSelect}
						bind:value={cpCharacter1}
						required
					>
						<option value="" selected disabled>Select a character</option>
						<option value="new">Create a new character</option>
						{#each characters as character}
							<option value={character.id}>{character.name}</option>
						{/each}
					</select>
				</label>

				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Your Character #2</span></div>
					<select
						name="cpCharacter2"
						placeholder="Select a character"
						class="select-bordered select w-full"
						autocomplete="name"
						onchange={handleChangeSelect}
						bind:value={cpCharacter2}
						required
					>
						<option value="" selected disabled>Select a character</option>
						<option value="new">Create a new character</option>
						{#each characters as character}
							<option value={character.id}>{character.name}</option>
						{/each}
					</select>
				</label>

				<p class="text-red-500">{errorText}</p>

				<div class="form-control mt-6">
					<button onclick={createCP} type="submit" class="btn w-full text-lg btn-primary"
						>Create</button
					>
				</div>
			</div>
		</div>
	</div>
</div>
