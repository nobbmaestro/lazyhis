#     __                      __  ___
#    / /   ____ _____  __  __/ / / (_)____
#   / /   / __ `/_  / / / / / /_/ / / ___/
#  / /___/ /_/ / / /_/ /_/ / __  / (__  )
# /_____/\__,_/ /___/\__, /_/ /_/_/____/
#                   /____/
#
# init.zsh - Shell initialization for lazyhis history integration
#
# Add the following to the end of your ~/.zshrc:
#     eval "$(lazyhis init zsh)"
#
# Inspired by atuin: https://github.com/atuinsh/atuin/blob/main/crates/atuin/src/shell/atuin.zsh
#
# ---
#
# MIT License (source: https://github.com/atuinsh/atuin/blob/main/LICENSE)
#
# Copyright (c) 2021 Ellie Huxtable
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

autoload -U add-zsh-hook

zmodload zsh/datetime 2>/dev/null

# If zsh-autosuggestions is installed, configure it to use LazyHis's search. If
# you'd like to override this, then add your config after the $(lazyhis init zsh)
# in your .zshrc
_zsh_autosuggest_strategy_lazyhis() {
	suggestion=$(lazyhis search --exit-code 0 --limit 1 -- "$@")
}

if [ -n "${ZSH_AUTOSUGGEST_STRATEGY:-}" ]; then
	ZSH_AUTOSUGGEST_STRATEGY=("lazyhis" "${ZSH_AUTOSUGGEST_STRATEGY[@]}")
else
	ZSH_AUTOSUGGEST_STRATEGY=("lazyhis")
fi

LAZYHIS_HISTORY_ID=""

_lazyhis_preexec() {
	local id
	id=$(lazyhis history add -- "$1")
	export LAZYHIS_HISTORY_ID="$id"
	__lazyhis_preexec_time=${EPOCHREALTIME-}
}

_lazyhis_precmd() {
	local EXIT="$?" __lazyhis_precmd_time=${EPOCHREALTIME-}

	[[ -z "${LAZYHIS_HISTORY_ID:-}" ]] && return

	local duration=""
	if [[ -n $__lazyhis_preexec_time && -n $__lazyhis_precmd_time ]]; then
		printf -v duration %.0f $(((__lazyhis_precmd_time - __lazyhis_preexec_time) * 1000))
	fi

	(
		lazyhis history edit \
			--exit-code $EXIT \
			${duration:+--duration=$duration} \
			-- $LAZYHIS_HISTORY_ID &
	) >/dev/null 2>&1

	export LAZYHIS_HISTORY_ID=""
}

_lazyhis_search() {
	emulate -L zsh
	zle -I

	local output
	output=$(lazyhis $* -- "$BUFFER" 3>&1 1>&2 2>&3)

	zle reset-prompt

	if [[ -n $output ]]; then
		RBUFFER=""
		LBUFFER=$output

		case $LBUFFER in
		__lazyhis_accept__:*)
			LBUFFER=${LBUFFER#__lazyhis_accept__:}
			zle accept-line
			;;
		__lazyhis_prefill__:*)
			LBUFFER=${LBUFFER#__lazyhis_prefill__:}
			;;
		esac
	fi
}
_lazyhis_search_vicmd() {
	_lazyhis_search
}
_lazyhis_search_viins() {
	_lazyhis_search
}

_lazyhis_up_search() {
	# Only trigger if the buffer is a single line
	if [[ ! $BUFFER == *$'\n'* ]]; then
		_lazyhis_search "$@"
	else
		zle up-line
	fi
}
_lazyhis_up_search_vicmd() {
	_lazyhis_up_search
}
_lazyhis_up_search_viins() {
	_lazyhis_up_search
}

add-zsh-hook preexec _lazyhis_preexec
add-zsh-hook precmd _lazyhis_precmd

zle -N lazyhis-search _lazyhis_search
zle -N lazyhis-search-vicmd _lazyhis_search_vicmd
zle -N lazyhis-search-viins _lazyhis_search_viins
zle -N lazyhis-up-search _lazyhis_up_search
zle -N lazyhis-up-search-vicmd _lazyhis_up_search_vicmd
zle -N lazyhis-up-search-viins _lazyhis_up_search_viins

bindkey -M emacs '^r' lazyhis-search
bindkey -M viins '^r' lazyhis-search-viins
bindkey -M vicmd '^r' lazyhis-search-vicmd

bindkey -M emacs '^[[A' lazyhis-up-search
bindkey -M vicmd '^[[A' lazyhis-up-search-vicmd
bindkey -M viins '^[[A' lazyhis-up-search-viins
bindkey -M emacs '^[OA' lazyhis-up-search
bindkey -M vicmd '^[OA' lazyhis-up-search-vicmd
bindkey -M viins '^[OA' lazyhis-up-search-viins
bindkey -M vicmd 'k' lazyhis-up-search-vicmd
